package file

// Input capabilities bind logical names to descriptor-acquired immutable bytes.
// They are not a filesystem sandbox or protection against privileged mounts.
import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/bearer/bearer/pkg/report/coverage"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"

	ignore "github.com/sabhiram/go-gitignore"
	"golang.org/x/sys/unix"
)

type InputIdentity struct {
	Device uint64 `json:"device"`
	Inode  uint64 `json:"inode"`
}
type InputExpectation struct {
	Root         string          `json:"root"`
	RootIdentity InputIdentity   `json:"root_identity"`
	Relative     string          `json:"relative"`
	Components   []InputIdentity `json:"components"`
	SHA256       string          `json:"sha256"`
}
type inputHandle struct {
	file     *os.File
	identity InputIdentity
	name     string
}

// Existence records contain metadata only: no processing outcome or file bytes.
// Repeated observations of the same logical marker retain one obligation.
type inputExistenceComponent struct {
	identity InputIdentity
	kind     uint32
}
type inputExistence struct {
	relative   string
	components []inputExistenceComponent
	present    bool
}
type InputScope struct {
	mu            sync.Mutex
	root          string
	logical       string
	single        string
	anchor        inputHandle
	inputs        []*BoundInput
	existence     map[string]inputExistence
	transfers     map[string]string
	contextInputs map[string]*BoundInput
	fault         error
	closed        bool
}
type BoundInput struct {
	scope           *InputScope
	relative        string
	handles         []inputHandle
	data            []byte
	info            os.FileInfo
	digest          string
	closed          bool
	validationReads uint64
	validationBytes uint64
	// Nonzero only for report-context samples; never a processing promise.
	sampleLimit int64
	context     bool
	size        int64
}
type InputReader struct {
	reader      *bytes.Reader
	input       *BoundInput
	closed      bool
	ownedScope  *InputScope
	reportInput *BoundInput
}

func identity(info os.FileInfo) (InputIdentity, error) {
	s, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return InputIdentity{}, errors.New("input identity unavailable")
	}
	return InputIdentity{uint64(s.Dev), uint64(s.Ino)}, nil
}
func handle(fd int, name string) (inputHandle, error) {
	f := os.NewFile(uintptr(fd), name)
	info, err := f.Stat()
	if err != nil {
		return inputHandle{}, errors.Join(err, f.Close())
	}
	id, err := identity(info)
	if err != nil {
		return inputHandle{}, errors.Join(err, f.Close())
	}
	return inputHandle{f, id, name}, nil
}
func OpenInputScope(root string) (*InputScope, error) {
	logical, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	// Resolve only the trusted root's parent, never a selected input leaf.
	info, err := os.Lstat(logical)
	if err != nil {
		return nil, err
	}
	single := ""
	anchor := logical
	if !info.IsDir() {
		if !info.Mode().IsRegular() {
			return nil, errors.New("non-regular input root")
		}
		single = filepath.Base(logical)
		anchor = filepath.Dir(logical)
	}
	anchor, err = filepath.EvalSymlinks(anchor)
	if err != nil {
		return nil, err
	}
	fd, err := unix.Open(anchor, unix.O_RDONLY|unix.O_DIRECTORY|unix.O_NOFOLLOW|unix.O_CLOEXEC, 0)
	if err != nil {
		return nil, err
	}
	h, err := handle(fd, anchor)
	if err != nil {
		return nil, err
	}
	return &InputScope{root: anchor, logical: logical, single: single, anchor: h}, nil
}
func (s *InputScope) latch(err error) error {
	if err != nil {
		s.fault = errors.Join(s.fault, err)
	}
	return err
}
func (s *InputScope) checkRoot() error {
	if s.closed {
		return s.latch(errors.New("input scope closed"))
	}
	fd, err := unix.Open(s.root, unix.O_RDONLY|unix.O_DIRECTORY|unix.O_NOFOLLOW|unix.O_CLOEXEC, 0)
	if err != nil {
		return s.latch(err)
	}
	h, err := handle(fd, s.root)
	if err != nil {
		return s.latch(err)
	}
	if h.identity != s.anchor.identity {
		err = errors.New("input root replaced")
	}
	return s.latch(errors.Join(err, h.file.Close()))
}
func (s *InputScope) selector(relative string) (string, error) {
	if strings.ContainsRune(relative, 0) || filepath.IsAbs(relative) {
		return "", errors.New("invalid input selector")
	}
	r := filepath.Clean(relative)
	if r == ".." || strings.HasPrefix(r, ".."+string(filepath.Separator)) {
		return "", errors.New("input outside target")
	}
	if s.single != "" {
		if r != "." {
			return "", errors.New("input outside single-file target")
		}
		return s.single, nil
	}
	if r == "." {
		return "", errors.New("input is not a regular file")
	}
	return r, nil
}
func (s *InputScope) Acquire(relative string, expected *InputExpectation) (*BoundInput, error) {
	return s.acquire(relative, expected, false)
}
func (s *InputScope) acquire(relative string, expected *InputExpectation, context bool) (*BoundInput, error) {
	return s.acquireLimited(relative, expected, context, 0)
}
func (s *InputScope) acquireLimited(relative string, expected *InputExpectation, context bool, sampleLimit int64) (*BoundInput, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.checkRoot(); err != nil {
		return nil, err
	}
	r, err := s.selector(relative)
	if context {
		r = filepath.Clean(relative)
		err = nil
		if strings.ContainsRune(relative, 0) || filepath.IsAbs(relative) || r == "." || r == ".." || strings.HasPrefix(r, ".."+string(filepath.Separator)) {
			err = errors.New("context outside target")
		}
	}
	if err != nil {
		return nil, s.latch(err)
	}
	if expected != nil {
		if expected.Root != s.root || expected.RootIdentity != s.anchor.identity || expected.Relative != relative || len(expected.SHA256) != 64 || len(expected.Components) == 0 {
			return nil, s.latch(errors.New("invalid input expectation"))
		}
	}
	b := &BoundInput{scope: s, relative: relative, sampleLimit: sampleLimit, context: context}
	parent := s.anchor.file
	parts := strings.Split(r, string(filepath.Separator))
	for i, part := range parts {
		flags := unix.O_RDONLY | unix.O_NOFOLLOW | unix.O_CLOEXEC | unix.O_NONBLOCK
		if i < len(parts)-1 {
			flags |= unix.O_DIRECTORY
		}
		fd, e := unix.Openat(int(parent.Fd()), part, flags, 0)
		if e != nil {
			err = e
			break
		}
		h, e := handle(fd, part)
		if e != nil {
			err = e
			break
		}
		b.handles = append(b.handles, h)
		parent = h.file
	}
	if err == nil {
		b.info, err = parent.Stat()
		if err == nil {
			b.size = b.info.Size()
		}
		if err == nil && !b.info.Mode().IsRegular() {
			err = errors.New("non-regular input")
		}
	}
	if err == nil && expected != nil {
		if len(expected.Components) != len(b.handles) {
			err = errors.New("input component expectation mismatch")
		} else {
			for i, h := range b.handles {
				if h.identity != expected.Components[i] {
					err = errors.New("input identity mismatch")
					break
				}
			}
		}
	}
	if err == nil {
		var reader io.Reader = parent
		if sampleLimit > 0 {
			reader = io.LimitReader(parent, sampleLimit)
		}
		b.data, err = io.ReadAll(reader)
	}
	if err == nil {
		sum := sha256.Sum256(b.data)
		b.digest = hex.EncodeToString(sum[:])
		if expected != nil && expected.SHA256 != b.digest {
			err = errors.New("input digest mismatch")
		}
	}
	if err == nil {
		err = b.validateLocked()
	}
	if err != nil {
		for _, h := range b.handles {
			err = errors.Join(err, h.file.Close())
		}
		return nil, s.latch(err)
	}
	s.inputs = append(s.inputs, b)
	return b, nil
}
func (b *BoundInput) validateIdentityLocked() error {
	s := b.scope
	if b.closed {
		return s.latch(errors.New("bound input closed"))
	}
	if err := s.checkRoot(); err != nil {
		return err
	}
	parent := s.anchor.file
	for i, h := range b.handles {
		flags := unix.O_RDONLY | unix.O_NOFOLLOW | unix.O_CLOEXEC | unix.O_NONBLOCK
		if i < len(b.handles)-1 {
			flags |= unix.O_DIRECTORY
		}
		fd, err := unix.Openat(int(parent.Fd()), h.name, flags, 0)
		if err != nil {
			return s.latch(err)
		}
		current, err := handle(fd, h.name)
		if err != nil {
			return s.latch(err)
		}
		if current.identity != h.identity {
			err = errors.New("input path replaced")
		}
		err = errors.Join(err, current.file.Close())
		if err != nil {
			return s.latch(err)
		}
		parent = h.file
	}
	return nil
}

// A full content check is performed at acquisition and explicit reconciliation,
// not for each small read of the already acquired immutable bytes.
func (b *BoundInput) validateContentLocked(f *os.File) error {
	b.validationReads++
	h := sha256.New()
	size := b.size
	info, err := f.Stat()
	if err != nil {
		return b.scope.latch(err)
	}
	if !info.Mode().IsRegular() || info.Size() != size {
		return b.scope.latch(errors.New("input metadata changed"))
	}
	readLimit, expectedSize := size+1, size
	if b.sampleLimit > 0 {
		expectedSize = min(size, b.sampleLimit)
		readLimit = expectedSize
	}
	n, err := io.CopyBuffer(h, io.NewSectionReader(f, 0, readLimit), make([]byte, 32*1024))
	b.validationBytes += uint64(n)
	if err == nil && (n != expectedSize || hex.EncodeToString(h.Sum(nil)) != b.digest) {
		err = errors.New("input contents changed")
	}
	return b.scope.latch(err)
}
func (b *BoundInput) validateLocked() error {
	if err := b.validateIdentityLocked(); err != nil {
		return err
	}
	return b.validateContentLocked(b.handles[len(b.handles)-1].file)
}

// Closing a per-file capability releases descriptors and bytes, not its final
// reconciliation obligation. Reacquire metadata/content through the held root,
// require every original identity, and close all temporary descriptors.
func (b *BoundInput) reconcileClosedLocked() (retErr error) {
	s := b.scope
	if err := s.checkRoot(); err != nil {
		return err
	}
	parent := s.anchor.file
	var held []*os.File
	defer func() {
		for i := len(held) - 1; i >= 0; i-- {
			retErr = errors.Join(retErr, held[i].Close())
		}
		s.latch(retErr)
	}()
	for i, expected := range b.handles {
		flags := unix.O_RDONLY | unix.O_NOFOLLOW | unix.O_CLOEXEC | unix.O_NONBLOCK
		if i < len(b.handles)-1 {
			flags |= unix.O_DIRECTORY
		}
		fd, err := unix.Openat(int(parent.Fd()), expected.name, flags, 0)
		if err != nil {
			return err
		}
		current, err := handle(fd, expected.name)
		if err != nil {
			return err
		}
		held = append(held, current.file)
		if current.identity != expected.identity {
			return errors.New("completed input path replaced")
		}
		parent = current.file
	}
	if len(held) == 0 {
		return errors.New("completed input provenance unavailable")
	}
	return b.validateContentLocked(parent)
}

func (b *BoundInput) Validate() error {
	b.scope.mu.Lock()
	defer b.scope.mu.Unlock()
	return b.validateLocked()
}
func (b *BoundInput) ReadBytes() ([]byte, error) {
	if err := b.Validate(); err != nil {
		return nil, err
	}
	return bytes.Clone(b.data), nil
}
func (b *BoundInput) Reader() (*InputReader, error) {
	if err := b.Validate(); err != nil {
		return nil, err
	}
	return &InputReader{reader: bytes.NewReader(b.data), input: b}, nil
}
func (r *InputReader) check() error {
	r.input.scope.mu.Lock()
	defer r.input.scope.mu.Unlock()
	if r.closed {
		return r.input.scope.latch(errors.New("input reader closed"))
	}
	return r.input.validateIdentityLocked()
}

func (r *InputReader) latchFailure(err error) error {
	if err != nil && err != io.EOF {
		r.input.scope.mu.Lock()
		r.input.scope.latch(err)
		r.input.scope.mu.Unlock()
	}
	return err
}

func (r *InputReader) Read(p []byte) (int, error) {
	if err := r.check(); err != nil {
		return 0, err
	}
	return r.reader.Read(p)
}
func (r *InputReader) ReadAt(p []byte, o int64) (int, error) {
	if err := r.check(); err != nil {
		return 0, err
	}
	n, err := r.reader.ReadAt(p, o)
	return n, r.latchFailure(err)
}
func (r *InputReader) Seek(o int64, w int) (int64, error) {
	if err := r.check(); err != nil {
		return 0, err
	}
	n, err := r.reader.Seek(o, w)
	return n, r.latchFailure(err)
}
func (r *InputReader) ReadByte() (byte, error) {
	if err := r.check(); err != nil {
		return 0, err
	}
	return r.reader.ReadByte()
}
func (r *InputReader) ReadRune() (rune, int, error) {
	if err := r.check(); err != nil {
		return 0, 0, err
	}
	return r.reader.ReadRune()
}
func (r *InputReader) UnreadByte() error {
	if err := r.check(); err != nil {
		return err
	}
	return r.latchFailure(r.reader.UnreadByte())
}
func (r *InputReader) UnreadRune() error {
	if err := r.check(); err != nil {
		return err
	}
	return r.latchFailure(r.reader.UnreadRune())
}
func (r *InputReader) WriteTo(w io.Writer) (int64, error) {
	if err := r.check(); err != nil {
		return 0, err
	}
	// Hide WriterTo so io.Copy uses the checked Read path on every refill.
	n, err := io.Copy(w, struct{ io.Reader }{r})
	if err != nil {
		r.input.scope.mu.Lock()
		r.input.scope.latch(err)
		r.input.scope.mu.Unlock()
	}
	return n, err
}
func (r *InputReader) Close() error {
	if r.closed {
		return nil
	}
	r.closed = true
	if r.reportInput != nil {
		return r.reportInput.Close()
	}
	if r.ownedScope != nil {
		return r.ownedScope.Close()
	}
	return nil
}
func (b *BoundInput) Digest() string    { return b.digest }
func (b *BoundInput) Info() os.FileInfo { return b.info }
func (b *BoundInput) Expectation() *InputExpectation {
	ids := make([]InputIdentity, len(b.handles))
	for i, h := range b.handles {
		ids[i] = h.identity
	}
	return &InputExpectation{b.scope.root, b.scope.anchor.identity, b.relative, ids, b.digest}
}
func (b *BoundInput) closeLocked() error {
	if b.closed {
		return nil
	}
	err := b.validateLocked()
	b.closed = true
	b.data = nil
	for i := len(b.handles) - 1; i >= 0; i-- {
		err = errors.Join(err, b.handles[i].file.Close())
	}
	return b.scope.latch(err)
}
func (b *BoundInput) Close() error {
	b.scope.mu.Lock()
	defer b.scope.mu.Unlock()
	return b.closeLocked()
}
func (b *BoundInput) Err() error { return b.scope.Err() }
func (s *InputScope) Err() error { s.mu.Lock(); defer s.mu.Unlock(); return s.fault }
func (s *InputScope) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return s.fault
	}
	for _, b := range s.inputs {
		if b.closed {
			b.reconcileClosedLocked()
		} else {
			b.closeLocked()
		}
	}
	s.reconcileExistenceLocked()
	s.checkRoot()
	s.closed = true
	s.latch(s.anchor.file.Close())
	return s.fault
}
func (s *InputScope) Validate() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.checkRoot(); err != nil {
		return err
	}
	for _, b := range s.inputs {
		if b.closed {
			b.reconcileClosedLocked()
		} else {
			b.validateLocked()
		}
	}
	s.reconcileExistenceLocked()
	return s.fault
}
func (p *Path) ReadBytes() (data []byte, retErr error) {
	if p.input != nil {
		return p.input.ReadBytes()
	}
	if p.scope != nil {
		r, err := filepath.Rel(p.scope.logicalRoot(), p.AbsolutePath)
		if err != nil {
			return nil, err
		}
		b, err := p.scope.acquire(r, nil, true)
		if err != nil {
			return nil, err
		}
		defer func() { retErr = errors.Join(retErr, b.Close()) }()
		return b.ReadBytes()
	}
	scope, err := OpenInputScope(p.AbsolutePath)
	if err != nil {
		return nil, err
	}
	defer func() { retErr = errors.Join(retErr, scope.Close()) }()
	b, err := scope.Acquire(".", nil)
	if err != nil {
		return nil, err
	}
	return b.ReadBytes()
}
func (p *Path) Reader() (*InputReader, error) {
	if p.input != nil {
		return p.input.Reader()
	}
	if p.scope != nil {
		relative, err := filepath.Rel(p.scope.logicalRoot(), p.AbsolutePath)
		if err != nil {
			return nil, err
		}
		b, err := p.scope.acquire(relative, nil, true)
		if err != nil {
			return nil, err
		}
		return b.Reader()
	}
	s, err := OpenInputScope(p.AbsolutePath)
	if err != nil {
		return nil, err
	}
	b, err := s.Acquire(".", nil)
	if err != nil {
		return nil, errors.Join(err, s.Close())
	}
	r, err := b.Reader()
	if err != nil {
		return nil, errors.Join(err, s.Close())
	}
	r.ownedScope = s
	return r, nil
}

// Optional .gitignore absence is not an acquisition failure. Any other failure
// is latched; the content itself is always acquired through the held directory.
func getScopedGitIgnore(s *InputScope) (*ignore.GitIgnore, error) {
	present := s.contextExists(filepath.Join(s.logicalRoot(), ".gitignore"))
	if err := s.Err(); err != nil {
		return nil, err
	}
	if !present {
		return nil, nil
	}
	b, err := s.acquire(".gitignore", nil, true)
	if err != nil {
		return nil, err
	}
	data, err := b.ReadBytes()
	err = errors.Join(err, b.Close())
	if err != nil {
		return nil, err
	}
	return ignore.CompileIgnoreLines(strings.Split(string(data), "\n")...), nil
}

// statRelative is metadata-only. Every component is still traversed relative to
// the held root without following links; it never authorizes a later path read.
func (s *InputScope) statRelative(relative string) (info os.FileInfo, retErr error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.checkRoot(); err != nil {
		return nil, err
	}
	r := filepath.Clean(relative)
	if s.single != "" {
		var err error
		r, err = s.selector(relative)
		if err != nil {
			return nil, s.latch(err)
		}
	}
	if strings.ContainsRune(r, 0) || filepath.IsAbs(r) || r == ".." || strings.HasPrefix(r, ".."+string(filepath.Separator)) {
		return nil, s.latch(errors.New("metadata outside target"))
	}
	if r == "." {
		return s.anchor.file.Stat()
	}
	parent := s.anchor.file
	var held []*os.File
	defer func() {
		for i := len(held) - 1; i >= 0; i-- {
			retErr = errors.Join(retErr, held[i].Close())
		}
		s.latch(retErr)
	}()
	parts := strings.Split(r, string(filepath.Separator))
	for i, part := range parts {
		flags := unix.O_RDONLY | unix.O_NOFOLLOW | unix.O_CLOEXEC | unix.O_NONBLOCK
		if i < len(parts)-1 {
			flags |= unix.O_DIRECTORY
		}
		fd, err := unix.Openat(int(parent.Fd()), part, flags, 0)
		if err != nil {
			return nil, err
		}
		parent = os.NewFile(uintptr(fd), part)
		held = append(held, parent)
	}
	return parent.Stat()
}

// Logical spelling is metadata; only the canonical held anchor performs I/O.
func (s *InputScope) logicalRoot() string {
	if s.single != "" {
		return filepath.Dir(s.logical)
	}
	return s.logical
}

// ReportReader reacquires only an input already bound in this same scope. It
// verifies the original component identities and digest before exposing bytes;
// it never promotes an arbitrary report filename into a new promised input.
func (s *InputScope) ReportReader(relative string) (*InputReader, error) {
	s.mu.Lock()
	if err := s.checkRoot(); err != nil {
		s.mu.Unlock()
		return nil, err
	}
	var expected *InputExpectation
	for _, b := range s.inputs {
		if b.relative == relative && b.sampleLimit == 0 {
			expected = b.Expectation()
			break
		}
	}
	if expected == nil {
		err := s.latch(errors.New("report input was not acquired during processing"))
		s.mu.Unlock()
		return nil, err
	}
	s.mu.Unlock()
	b, err := s.Acquire(relative, expected)
	if err != nil {
		return nil, err
	}
	r, err := b.Reader()
	if err != nil {
		return nil, errors.Join(err, b.Close())
	}
	// This reader owns just the reacquired handle; original reconciliation entries
	// stay in the scope until final report acceptance.
	r.reportInput = b
	return r, nil
}

// ReportSample reads a bounded language sample. Already processed inputs are
// reacquired against their original full-content expectation. Other context is
// root-confined and tracked only as sampled report context, never as processing.
// The sample digest, file size and component identities survive reader closure
// until final report acceptance. Unsampled context bytes are not coverage.
func (s *InputScope) ReportSample(relative string, limit int) (sample []byte, info os.FileInfo, retErr error) {
	if limit <= 0 || limit > 16*1024 {
		s.mu.Lock()
		defer s.mu.Unlock()
		return nil, nil, s.latch(errors.New("invalid report sample limit"))
	}
	s.mu.Lock()
	if err := s.checkRoot(); err != nil {
		s.mu.Unlock()
		return nil, nil, err
	}
	var expected *InputExpectation
	for _, b := range s.inputs {
		if b.relative == relative && b.sampleLimit == 0 {
			expected = b.Expectation()
			break
		}
	}
	s.mu.Unlock()
	var b *BoundInput
	var err error
	if expected != nil {
		b, err = s.Acquire(relative, expected)
	} else {
		b, err = s.acquireLimited(relative, nil, false, int64(limit))
	}
	if err != nil {
		return nil, nil, err
	}
	defer func() { retErr = errors.Join(retErr, b.Close()) }()
	reader, err := b.Reader()
	if err != nil {
		return nil, nil, err
	}
	defer func() { retErr = errors.Join(retErr, reader.Close()) }()
	sample, err = io.ReadAll(io.LimitReader(reader, int64(limit)))
	return sample, b.Info(), err
}

// contextExists uses the same acquired root as detector reads. Normal absence
// is an observation, not an acquisition fault; all other failures are terminal.
func (s *InputScope) contextExists(absolute string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.checkRoot(); err != nil {
		return false
	}
	relative, err := filepath.Rel(s.logicalRoot(), absolute)
	if err != nil {
		s.latch(err)
		return false
	}
	observation, err := s.observeExistenceLocked(relative)
	if err != nil {
		s.latch(err)
		return false
	}
	if previous, ok := s.existence[relative]; ok && !sameExistence(previous, observation) {
		s.latch(errors.New("context existence changed"))
		return false
	}
	if s.existence == nil {
		s.existence = make(map[string]inputExistence)
	}
	s.existence[relative] = observation
	return observation.present
}

// Every pathname lookup is relative to an owned descriptor. The last component
// is inspected with NOFOLLOW rather than opened for content; legitimate metadata
// checks do not require read access to the file or read its contents.
func (s *InputScope) observeExistenceLocked(relative string) (observation inputExistence, retErr error) {
	observation.relative = relative
	if strings.ContainsRune(relative, 0) || filepath.IsAbs(relative) {
		return observation, errors.New("invalid context selector")
	}
	r := filepath.Clean(relative)
	if r == ".." || strings.HasPrefix(r, ".."+string(filepath.Separator)) {
		return observation, errors.New("context outside target")
	}
	parent := s.anchor.file
	var held []*os.File
	defer func() {
		for i := len(held) - 1; i >= 0; i-- {
			retErr = errors.Join(retErr, held[i].Close())
		}
	}()
	parts := strings.Split(r, string(filepath.Separator))
	for i, part := range parts {
		var stat unix.Stat_t
		err := unix.Fstatat(int(parent.Fd()), part, &stat, unix.AT_SYMLINK_NOFOLLOW)
		if errors.Is(err, unix.ENOENT) {
			return observation, nil
		}
		if err != nil {
			return observation, err
		}
		kind := uint32(stat.Mode) & unix.S_IFMT
		if kind == unix.S_IFLNK {
			return observation, errors.New("context symlink rejected")
		}
		component := inputExistenceComponent{identity: InputIdentity{uint64(stat.Dev), uint64(stat.Ino)}, kind: kind}
		observation.components = append(observation.components, component)
		if i == len(parts)-1 {
			observation.present = true
			return observation, nil
		}
		if kind != unix.S_IFDIR {
			return observation, errors.New("context component is not a directory")
		}
		fd, err := unix.Openat(int(parent.Fd()), part, unix.O_RDONLY|unix.O_DIRECTORY|unix.O_NOFOLLOW|unix.O_CLOEXEC, 0)
		if err != nil {
			return observation, err
		}
		current, err := handle(fd, part)
		if err != nil {
			return observation, err
		}
		held = append(held, current.file)
		if current.identity != component.identity {
			return observation, errors.New("context component replaced")
		}
		parent = current.file
	}
	return observation, errors.New("context observation unavailable")
}

func sameExistence(a, b inputExistence) bool {
	if a.relative != b.relative || a.present != b.present || len(a.components) != len(b.components) {
		return false
	}
	for i := range a.components {
		if a.components[i] != b.components[i] {
			return false
		}
	}
	return true
}

func (s *InputScope) reconcileExistenceLocked() {
	if len(s.existence) == 0 {
		return
	}
	if err := s.checkRoot(); err != nil {
		return
	}
	for relative, expected := range s.existence {
		observed, err := s.observeExistenceLocked(relative)
		if err == nil && !sameExistence(expected, observed) {
			err = errors.New("context existence changed")
		}
		s.latch(err)
	}
}

// ContextReceipt is an internal same-run handoff, never processing coverage.
// Even an empty observation list is explicit; a missing receipt is not empty.
const ContextReceiptType = "input_context"
const contextReceiptSchema = "bearer-context-v2"
const legacyContextReceiptSchema = "bearer-context-v1"

type ContextComponent struct {
	Identity InputIdentity `json:"identity"`
	Kind     uint32        `json:"kind"`
}
type ContextObservation struct {
	Relative   string             `json:"relative"`
	Components []ContextComponent `json:"components"`
	Present    bool               `json:"present"`
}

// ContextRead preserves the acquisition record, not a new observation at import.
// A zero sample limit means all bytes; bounded samples retain their original size.
type ContextRead struct {
	Input       InputExpectation `json:"input"`
	Size        int64            `json:"size"`
	SampleLimit int64            `json:"sample_limit"`
}
type ContextReceipt struct {
	Type         string               `json:"type"`
	Schema       string               `json:"schema"`
	Input        InputExpectation     `json:"input"`
	Observations []ContextObservation `json:"observations"`
	Reads        []ContextRead        `json:"reads"`
}

// Legacy cache records remain readable, but are never sufficient fresh handoffs.
func (r ContextReceipt) MarshalJSON() ([]byte, error) {
	if r.Schema == legacyContextReceiptSchema && r.Reads == nil {
		return json.Marshal(struct {
			Type         string               `json:"type"`
			Schema       string               `json:"schema"`
			Input        InputExpectation     `json:"input"`
			Observations []ContextObservation `json:"observations"`
		}{r.Type, r.Schema, r.Input, r.Observations})
	}
	type plain ContextReceipt
	return json.Marshal(plain(r))
}

func contextObject(raw json.RawMessage, keys ...string) (map[string]json.RawMessage, error) {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(raw, &obj); err != nil {
		return nil, err
	}
	if obj == nil || len(obj) != len(keys) {
		return nil, errors.New("invalid context receipt fields")
	}
	for _, key := range keys {
		v, ok := obj[key]
		if !ok || bytes.Equal(bytes.TrimSpace(v), []byte("null")) {
			return nil, errors.New("missing/null context receipt field")
		}
	}
	return obj, nil
}
func contextArray(raw json.RawMessage) ([]json.RawMessage, error) {
	var values []json.RawMessage
	if err := json.Unmarshal(raw, &values); err != nil {
		return nil, err
	}
	if values == nil {
		return nil, errors.New("null context receipt array")
	}
	return values, nil
}
func contextIdentity(raw json.RawMessage) error {
	_, err := contextObject(raw, "device", "inode")
	return err
}
func contextExpectation(raw json.RawMessage) error {
	fields, err := contextObject(raw, "root", "root_identity", "relative", "components", "sha256")
	if err != nil {
		return err
	}
	if err = contextIdentity(fields["root_identity"]); err != nil {
		return err
	}
	cs, err := contextArray(fields["components"])
	if err != nil {
		return err
	}
	for _, c := range cs {
		if err = contextIdentity(c); err != nil {
			return err
		}
	}
	return nil
}
func contextReadKey(relative string, limit int64) string {
	return relative + "\x00" + strconv.FormatInt(limit, 10)
}
func DecodeContextReceipt(raw []byte) (*ContextReceipt, error) {
	if err := coverage.CheckJSON(raw); err != nil {
		return nil, err
	}
	var header struct {
		Schema string `json:"schema"`
	}
	if err := json.Unmarshal(raw, &header); err != nil {
		return nil, err
	}
	keys := []string{"type", "schema", "input", "observations"}
	if header.Schema != legacyContextReceiptSchema {
		keys = append(keys, "reads")
	}
	obj, err := contextObject(raw, keys...)
	if err != nil {
		return nil, err
	}
	if err = contextExpectation(obj["input"]); err != nil {
		return nil, err
	}
	if header.Schema != legacyContextReceiptSchema {
		reads, e := contextArray(obj["reads"])
		if e != nil {
			return nil, e
		}
		for _, rawRead := range reads {
			fields, e := contextObject(rawRead, "input", "size", "sample_limit")
			if e != nil {
				return nil, e
			}
			if e = contextExpectation(fields["input"]); e != nil {
				return nil, e
			}
		}
	}
	observations, err := contextArray(obj["observations"])
	if err != nil {
		return nil, err
	}
	for _, rawObservation := range observations {
		observation, e := contextObject(rawObservation, "relative", "components", "present")
		if e != nil {
			return nil, e
		}
		cs, e := contextArray(observation["components"])
		if e != nil {
			return nil, e
		}
		for _, c := range cs {
			fields, e := contextObject(c, "identity", "kind")
			if e != nil {
				return nil, e
			}
			if e = contextIdentity(fields["identity"]); e != nil {
				return nil, e
			}
		}
	}
	var result ContextReceipt
	if err = json.Unmarshal(raw, &result); err != nil {
		return nil, err
	}
	if err = result.validate(); err != nil {
		return nil, err
	}
	return &result, nil
}
func (r *ContextReceipt) validate() error {
	if r == nil || r.Type != ContextReceiptType || (r.Schema != contextReceiptSchema && r.Schema != legacyContextReceiptSchema) || r.Observations == nil || len(r.Input.Components) == 0 || !filepath.IsAbs(r.Input.Root) {
		return errors.New("invalid context receipt")
	}
	digest, e := hex.DecodeString(r.Input.SHA256)
	if e != nil || len(digest) != sha256.Size {
		return errors.New("invalid context input digest")
	}
	safe := func(name string) bool {
		return name != "" && !strings.ContainsRune(name, 0) && !filepath.IsAbs(name) && filepath.Clean(name) == name && name != ".." && !strings.HasPrefix(name, ".."+string(filepath.Separator))
	}
	if !safe(r.Input.Relative) {
		return errors.New("invalid context input selector")
	}
	if r.Schema == contextReceiptSchema && r.Reads == nil {
		return errors.New("missing context read provenance")
	}
	if r.Schema == legacyContextReceiptSchema && r.Reads != nil {
		return errors.New("unexpected legacy context reads")
	}
	reads := map[string]bool{}
	for _, c := range r.Reads {
		e := c.Input
		digest, err := hex.DecodeString(e.SHA256)
		key := contextReadKey(e.Relative, c.SampleLimit)
		if !safe(e.Relative) || e.Relative == "." || e.Root != r.Input.Root || e.RootIdentity != r.Input.RootIdentity || len(e.Components) != len(strings.Split(e.Relative, string(filepath.Separator))) || err != nil || len(digest) != sha256.Size || c.Size < 0 || c.Size == int64(^uint64(0)>>1) || c.SampleLimit < 0 || c.SampleLimit > 16*1024 || reads[key] {
			return errors.New("invalid/duplicate context read provenance")
		}
		reads[key] = true
	}
	seen := map[string]bool{}
	for _, o := range r.Observations {
		if !safe(o.Relative) || seen[o.Relative] || o.Components == nil {
			return errors.New("invalid/duplicate context observation")
		}
		seen[o.Relative] = true
		count := len(strings.Split(o.Relative, string(filepath.Separator)))
		if (o.Present && len(o.Components) != count) || (!o.Present && len(o.Components) >= count) {
			return errors.New("invalid context observation shape")
		}
		for i, c := range o.Components {
			switch c.Kind {
			case unix.S_IFREG, unix.S_IFDIR, unix.S_IFIFO, unix.S_IFSOCK, unix.S_IFCHR, unix.S_IFBLK:
			default:
				return errors.New("invalid context component type")
			}
			if (i < len(o.Components)-1 || !o.Present) && c.Kind != unix.S_IFDIR {
				return errors.New("non-directory context prefix")
			}
		}
	}
	return nil
}
func sameExpectation(a, b *InputExpectation) bool {
	if a == nil || b == nil || a.Root != b.Root || a.RootIdentity != b.RootIdentity || a.Relative != b.Relative || a.SHA256 != b.SHA256 || len(a.Components) != len(b.Components) {
		return false
	}
	for i := range a.Components {
		if a.Components[i] != b.Components[i] {
			return false
		}
	}
	return true
}
func (s *InputScope) matchesInputLocked(expected *InputExpectation) bool {
	if expected == nil || expected.Root != s.root || expected.RootIdentity != s.anchor.identity {
		return false
	}
	for _, b := range s.inputs {
		if !b.context && b.sampleLimit == 0 && sameExpectation(b.Expectation(), expected) {
			return true
		}
	}
	return false
}
func (s *InputScope) ContextReceipt(expected *InputExpectation) (*ContextReceipt, error) {
	if err := s.Validate(); err != nil {
		return nil, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return nil, s.latch(errors.New("input scope closed"))
	}
	if !s.matchesInputLocked(expected) {
		return nil, s.latch(errors.New("context receipt input mismatch"))
	}
	result := &ContextReceipt{Type: ContextReceiptType, Schema: contextReceiptSchema, Input: *expected, Observations: []ContextObservation{}, Reads: []ContextRead{}}
	result.Input.Components = append([]InputIdentity{}, expected.Components...)
	names := make([]string, 0, len(s.existence))
	for name := range s.existence {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		o := s.existence[name]
		out := ContextObservation{Relative: name, Present: o.present, Components: []ContextComponent{}}
		for _, c := range o.components {
			out.Components = append(out.Components, ContextComponent{c.identity, c.kind})
		}
		result.Observations = append(result.Observations, out)
	}
	// Closed capabilities retain their immutable acquisition records. No pathname
	// lookup here can replace the bytes/identities that extraction actually used.
	reads := map[string]ContextRead{}
	for _, b := range s.inputs {
		if !b.context {
			continue
		}
		c := ContextRead{Input: *b.Expectation(), Size: b.size, SampleLimit: b.sampleLimit}
		key := contextReadKey(c.Input.Relative, c.SampleLimit)
		if old, ok := reads[key]; ok && (!sameExpectation(&old.Input, &c.Input) || old.Size != c.Size) {
			return nil, s.latch(errors.New("conflicting acquired context"))
		}
		reads[key] = c
	}
	names = names[:0]
	for name := range reads {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		result.Reads = append(result.Reads, reads[name])
	}
	if err := result.validate(); err != nil {
		return nil, s.latch(err)
	}
	return result, nil
}
func (s *InputScope) ImportContextReceipt(receipt *ContextReceipt, expected *InputExpectation) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.checkRoot(); err != nil {
		return err
	}
	if err := receipt.validate(); err != nil {
		return s.latch(err)
	}
	if receipt.Schema != contextReceiptSchema {
		return s.latch(errors.New("fresh context read provenance unavailable"))
	}
	if (expected != nil && !sameExpectation(expected, &receipt.Input)) || !s.matchesInputLocked(&receipt.Input) {
		return s.latch(errors.New("context handoff input/root mismatch"))
	}
	encoded, err := json.Marshal(receipt)
	if err != nil {
		return s.latch(err)
	}
	key := receipt.Input.Relative
	if old, ok := s.transfers[key]; ok && old != string(encoded) {
		return s.latch(errors.New("conflicting context handoff"))
	}
	if s.existence == nil {
		s.existence = map[string]inputExistence{}
	}
	for _, o := range receipt.Observations {
		retained := inputExistence{relative: o.Relative, present: o.Present}
		for _, c := range o.Components {
			retained.components = append(retained.components, inputExistenceComponent{c.Identity, c.Kind})
		}
		if old, ok := s.existence[o.Relative]; ok && !sameExistence(old, retained) {
			return s.latch(errors.New("conflicting context observation"))
		}
		// Reconcile the transferred observation; never replace it with today's state.
		actual, e := s.observeExistenceLocked(o.Relative)
		if e != nil {
			return s.latch(e)
		}
		if !sameExistence(actual, retained) {
			return s.latch(errors.New("context handoff changed"))
		}
		s.existence[o.Relative] = retained
	}
	if s.contextInputs == nil {
		s.contextInputs = map[string]*BoundInput{}
	}
	for _, c := range receipt.Reads {
		key := contextReadKey(c.Input.Relative, c.SampleLimit)
		retained, ok := s.contextInputs[key]
		if ok {
			if !sameExpectation(retained.Expectation(), &c.Input) || retained.size != c.Size {
				return s.latch(errors.New("conflicting context read handoff"))
			}
		} else {
			// A closed reconciliation-only record owns no file descriptor or bytes.
			// Reuse the existing held-root NOFOLLOW reconciliation and final-close path.
			retained = &BoundInput{scope: s, relative: c.Input.Relative, digest: c.Input.SHA256, size: c.Size, sampleLimit: c.SampleLimit, context: true, closed: true}
			for i, name := range strings.Split(c.Input.Relative, string(filepath.Separator)) {
				retained.handles = append(retained.handles, inputHandle{identity: c.Input.Components[i], name: name})
			}
			s.inputs = append(s.inputs, retained)
			s.contextInputs[key] = retained
		}
		if err := retained.reconcileClosedLocked(); err != nil {
			return s.latch(err)
		}
	}
	if s.transfers == nil {
		s.transfers = map[string]string{}
	}
	s.transfers[key] = string(encoded)
	return s.fault
}
func (s *InputScope) RequireContextReceipts(names []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.checkRoot(); err != nil {
		return err
	}
	if len(names) != len(s.transfers) {
		return s.latch(errors.New("missing/unexpected context handoff"))
	}
	seen := map[string]bool{}
	for _, name := range names {
		if seen[name] || s.transfers[name] == "" {
			return s.latch(errors.New("missing/duplicate context handoff"))
		}
		seen[name] = true
	}
	return s.fault
}

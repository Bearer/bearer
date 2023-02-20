package proto_test

import (
	"context"
	"testing"

	"github.com/bearer/bearer/pkg/parser/sitter/proto"
	"github.com/stretchr/testify/assert"

	sitter "github.com/smacker/go-tree-sitter"
)

type NodeContent struct {
	Type    string
	Content string
}

func TestGrammar(t *testing.T) {
	input := []byte(`syntax = "proto3";

	import "google/protobuf/timestamp.proto";

	option csharp_namespace = "Import.Api";

	package member;

	service Import {
		rpc ImportMembers (MembersBatch) returns (ImportReply);
		rpc ImportAddresses (AddressesBatch) returns (ImportReply);
	}

	message Member {
		string member_id = 1;
		string firstname = 2;
		string lastname = 3;
		string gender = 4;
		string external_user_id = 5;
		google.protobuf.Timestamp birthdate = 6;
		google.protobuf.Timestamp modification_date = 7;
		string shop_name = 8;
		string status = 9;
		google.protobuf.Timestamp registration_date = 10;
		string email = 11;
		string password = 12;
	}`)
	rootNode, err := sitter.ParseCtx(context.Background(), input, proto.GetLanguage())
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(
		t,
		`(source_file (syntax) (import path: (string)) (option (identifier) (constant (string))) (package (full_ident (identifier))) (service (service_name (identifier)) (rpc (rpc_name (identifier)) (message_or_enum_type (identifier)) (message_or_enum_type (identifier))) (rpc (rpc_name (identifier)) (message_or_enum_type (identifier)) (message_or_enum_type (identifier)))) (message (message_name (identifier)) (message_body (field (type) (identifier) (field_number (int_lit (decimal_lit)))) (field (type) (identifier) (field_number (int_lit (decimal_lit)))) (field (type) (identifier) (field_number (int_lit (decimal_lit)))) (field (type) (identifier) (field_number (int_lit (decimal_lit)))) (field (type) (identifier) (field_number (int_lit (decimal_lit)))) (field (type (message_or_enum_type (identifier) (identifier) (identifier))) (identifier) (field_number (int_lit (decimal_lit)))) (field (type (message_or_enum_type (identifier) (identifier) (identifier))) (identifier) (field_number (int_lit (decimal_lit)))) (field (type) (identifier) (field_number (int_lit (decimal_lit)))) (field (type) (identifier) (field_number (int_lit (decimal_lit)))) (field (type (message_or_enum_type (identifier) (identifier) (identifier))) (identifier) (field_number (int_lit (decimal_lit)))) (field (type) (identifier) (field_number (int_lit (decimal_lit)))) (field (type) (identifier) (field_number (int_lit (decimal_lit)))))))`,
		rootNode.String(),
	)
}

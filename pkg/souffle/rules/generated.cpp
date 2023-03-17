#define SOUFFLE_GENERATOR_VERSION "2.4-1-g15b114abf"
#include "souffle/CompiledSouffle.h"
#include "souffle/SignalHandler.h"
#include "souffle/SouffleInterface.h"
#include "souffle/datastructure/BTree.h"
#include "souffle/io/IOSystem.h"
#include <any>
namespace functors {
extern "C" {
}
} //namespace functors
namespace souffle::t_btree_ii__0_1__11 {
using namespace souffle;
struct Type {
static constexpr Relation::arity_type Arity = 2;
using t_tuple = Tuple<RamDomain, 2>;
struct t_comparator_0{
 int operator()(const t_tuple& a, const t_tuple& b) const {
  return (ramBitCast<RamSigned>(a[0]) < ramBitCast<RamSigned>(b[0])) ? -1 : (ramBitCast<RamSigned>(a[0]) > ramBitCast<RamSigned>(b[0])) ? 1 :((ramBitCast<RamSigned>(a[1]) < ramBitCast<RamSigned>(b[1])) ? -1 : (ramBitCast<RamSigned>(a[1]) > ramBitCast<RamSigned>(b[1])) ? 1 :(0));
 }
bool less(const t_tuple& a, const t_tuple& b) const {
  return (ramBitCast<RamSigned>(a[0]) < ramBitCast<RamSigned>(b[0]))|| ((ramBitCast<RamSigned>(a[0]) == ramBitCast<RamSigned>(b[0])) && ((ramBitCast<RamSigned>(a[1]) < ramBitCast<RamSigned>(b[1]))));
 }
bool equal(const t_tuple& a, const t_tuple& b) const {
return (ramBitCast<RamSigned>(a[0]) == ramBitCast<RamSigned>(b[0]))&&(ramBitCast<RamSigned>(a[1]) == ramBitCast<RamSigned>(b[1]));
 }
};
using t_ind_0 = btree_set<t_tuple,t_comparator_0>;
t_ind_0 ind_0;
using iterator = t_ind_0::iterator;
struct context {
t_ind_0::operation_hints hints_0_lower;
t_ind_0::operation_hints hints_0_upper;
};
context createContext() { return context(); }
bool insert(const t_tuple& t);
bool insert(const t_tuple& t, context& h);
bool insert(const RamDomain* ramDomain);
bool insert(RamDomain a0,RamDomain a1);
bool contains(const t_tuple& t, context& h) const;
bool contains(const t_tuple& t) const;
std::size_t size() const;
iterator find(const t_tuple& t, context& h) const;
iterator find(const t_tuple& t) const;
range<iterator> lowerUpperRange_00(const t_tuple& /* lower */, const t_tuple& /* upper */, context& /* h */) const;
range<iterator> lowerUpperRange_00(const t_tuple& /* lower */, const t_tuple& /* upper */) const;
range<t_ind_0::iterator> lowerUpperRange_11(const t_tuple& lower, const t_tuple& upper, context& h) const;
range<t_ind_0::iterator> lowerUpperRange_11(const t_tuple& lower, const t_tuple& upper) const;
bool empty() const;
std::vector<range<iterator>> partition() const;
void purge();
iterator begin() const;
iterator end() const;
void printStatistics(std::ostream& o) const;
};
} // namespace souffle::t_btree_ii__0_1__11 
namespace souffle::t_btree_ii__0_1__11 {
using namespace souffle;
using t_ind_0 = Type::t_ind_0;
using iterator = Type::iterator;
using context = Type::context;
bool Type::insert(const t_tuple& t) {
context h;
return insert(t, h);
}
bool Type::insert(const t_tuple& t, context& h) {
if (ind_0.insert(t, h.hints_0_lower)) {
return true;
} else return false;
}
bool Type::insert(const RamDomain* ramDomain) {
RamDomain data[2];
std::copy(ramDomain, ramDomain + 2, data);
const t_tuple& tuple = reinterpret_cast<const t_tuple&>(data);
context h;
return insert(tuple, h);
}
bool Type::insert(RamDomain a0,RamDomain a1) {
RamDomain data[2] = {a0,a1};
return insert(data);
}
bool Type::contains(const t_tuple& t, context& h) const {
return ind_0.contains(t, h.hints_0_lower);
}
bool Type::contains(const t_tuple& t) const {
context h;
return contains(t, h);
}
std::size_t Type::size() const {
return ind_0.size();
}
iterator Type::find(const t_tuple& t, context& h) const {
return ind_0.find(t, h.hints_0_lower);
}
iterator Type::find(const t_tuple& t) const {
context h;
return find(t, h);
}
range<iterator> Type::lowerUpperRange_00(const t_tuple& /* lower */, const t_tuple& /* upper */, context& /* h */) const {
return range<iterator>(ind_0.begin(),ind_0.end());
}
range<iterator> Type::lowerUpperRange_00(const t_tuple& /* lower */, const t_tuple& /* upper */) const {
return range<iterator>(ind_0.begin(),ind_0.end());
}
range<t_ind_0::iterator> Type::lowerUpperRange_11(const t_tuple& lower, const t_tuple& upper, context& h) const {
t_comparator_0 comparator;
int cmp = comparator(lower, upper);
if (cmp == 0) {
    auto pos = ind_0.find(lower, h.hints_0_lower);
    auto fin = ind_0.end();
    if (pos != fin) {fin = pos; ++fin;}
    return make_range(pos, fin);
}
if (cmp > 0) {
    return make_range(ind_0.end(), ind_0.end());
}
return make_range(ind_0.lower_bound(lower, h.hints_0_lower), ind_0.upper_bound(upper, h.hints_0_upper));
}
range<t_ind_0::iterator> Type::lowerUpperRange_11(const t_tuple& lower, const t_tuple& upper) const {
context h;
return lowerUpperRange_11(lower,upper,h);
}
bool Type::empty() const {
return ind_0.empty();
}
std::vector<range<iterator>> Type::partition() const {
return ind_0.getChunks(400);
}
void Type::purge() {
ind_0.clear();
}
iterator Type::begin() const {
return ind_0.begin();
}
iterator Type::end() const {
return ind_0.end();
}
void Type::printStatistics(std::ostream& o) const {
o << " arity 2 direct b-tree index 0 lex-order [0,1]\n";
ind_0.printStats(o);
}
} // namespace souffle::t_btree_ii__0_1__11 
namespace souffle::t_btree_iii__0_2_1__101__111 {
using namespace souffle;
struct Type {
static constexpr Relation::arity_type Arity = 3;
using t_tuple = Tuple<RamDomain, 3>;
struct t_comparator_0{
 int operator()(const t_tuple& a, const t_tuple& b) const {
  return (ramBitCast<RamSigned>(a[0]) < ramBitCast<RamSigned>(b[0])) ? -1 : (ramBitCast<RamSigned>(a[0]) > ramBitCast<RamSigned>(b[0])) ? 1 :((ramBitCast<RamSigned>(a[2]) < ramBitCast<RamSigned>(b[2])) ? -1 : (ramBitCast<RamSigned>(a[2]) > ramBitCast<RamSigned>(b[2])) ? 1 :((ramBitCast<RamSigned>(a[1]) < ramBitCast<RamSigned>(b[1])) ? -1 : (ramBitCast<RamSigned>(a[1]) > ramBitCast<RamSigned>(b[1])) ? 1 :(0)));
 }
bool less(const t_tuple& a, const t_tuple& b) const {
  return (ramBitCast<RamSigned>(a[0]) < ramBitCast<RamSigned>(b[0]))|| ((ramBitCast<RamSigned>(a[0]) == ramBitCast<RamSigned>(b[0])) && ((ramBitCast<RamSigned>(a[2]) < ramBitCast<RamSigned>(b[2]))|| ((ramBitCast<RamSigned>(a[2]) == ramBitCast<RamSigned>(b[2])) && ((ramBitCast<RamSigned>(a[1]) < ramBitCast<RamSigned>(b[1]))))));
 }
bool equal(const t_tuple& a, const t_tuple& b) const {
return (ramBitCast<RamSigned>(a[0]) == ramBitCast<RamSigned>(b[0]))&&(ramBitCast<RamSigned>(a[2]) == ramBitCast<RamSigned>(b[2]))&&(ramBitCast<RamSigned>(a[1]) == ramBitCast<RamSigned>(b[1]));
 }
};
using t_ind_0 = btree_set<t_tuple,t_comparator_0>;
t_ind_0 ind_0;
using iterator = t_ind_0::iterator;
struct context {
t_ind_0::operation_hints hints_0_lower;
t_ind_0::operation_hints hints_0_upper;
};
context createContext() { return context(); }
bool insert(const t_tuple& t);
bool insert(const t_tuple& t, context& h);
bool insert(const RamDomain* ramDomain);
bool insert(RamDomain a0,RamDomain a1,RamDomain a2);
bool contains(const t_tuple& t, context& h) const;
bool contains(const t_tuple& t) const;
std::size_t size() const;
iterator find(const t_tuple& t, context& h) const;
iterator find(const t_tuple& t) const;
range<iterator> lowerUpperRange_000(const t_tuple& /* lower */, const t_tuple& /* upper */, context& /* h */) const;
range<iterator> lowerUpperRange_000(const t_tuple& /* lower */, const t_tuple& /* upper */) const;
range<t_ind_0::iterator> lowerUpperRange_101(const t_tuple& lower, const t_tuple& upper, context& h) const;
range<t_ind_0::iterator> lowerUpperRange_101(const t_tuple& lower, const t_tuple& upper) const;
range<t_ind_0::iterator> lowerUpperRange_111(const t_tuple& lower, const t_tuple& upper, context& h) const;
range<t_ind_0::iterator> lowerUpperRange_111(const t_tuple& lower, const t_tuple& upper) const;
bool empty() const;
std::vector<range<iterator>> partition() const;
void purge();
iterator begin() const;
iterator end() const;
void printStatistics(std::ostream& o) const;
};
} // namespace souffle::t_btree_iii__0_2_1__101__111 
namespace souffle::t_btree_iii__0_2_1__101__111 {
using namespace souffle;
using t_ind_0 = Type::t_ind_0;
using iterator = Type::iterator;
using context = Type::context;
bool Type::insert(const t_tuple& t) {
context h;
return insert(t, h);
}
bool Type::insert(const t_tuple& t, context& h) {
if (ind_0.insert(t, h.hints_0_lower)) {
return true;
} else return false;
}
bool Type::insert(const RamDomain* ramDomain) {
RamDomain data[3];
std::copy(ramDomain, ramDomain + 3, data);
const t_tuple& tuple = reinterpret_cast<const t_tuple&>(data);
context h;
return insert(tuple, h);
}
bool Type::insert(RamDomain a0,RamDomain a1,RamDomain a2) {
RamDomain data[3] = {a0,a1,a2};
return insert(data);
}
bool Type::contains(const t_tuple& t, context& h) const {
return ind_0.contains(t, h.hints_0_lower);
}
bool Type::contains(const t_tuple& t) const {
context h;
return contains(t, h);
}
std::size_t Type::size() const {
return ind_0.size();
}
iterator Type::find(const t_tuple& t, context& h) const {
return ind_0.find(t, h.hints_0_lower);
}
iterator Type::find(const t_tuple& t) const {
context h;
return find(t, h);
}
range<iterator> Type::lowerUpperRange_000(const t_tuple& /* lower */, const t_tuple& /* upper */, context& /* h */) const {
return range<iterator>(ind_0.begin(),ind_0.end());
}
range<iterator> Type::lowerUpperRange_000(const t_tuple& /* lower */, const t_tuple& /* upper */) const {
return range<iterator>(ind_0.begin(),ind_0.end());
}
range<t_ind_0::iterator> Type::lowerUpperRange_101(const t_tuple& lower, const t_tuple& upper, context& h) const {
t_comparator_0 comparator;
int cmp = comparator(lower, upper);
if (cmp > 0) {
    return make_range(ind_0.end(), ind_0.end());
}
return make_range(ind_0.lower_bound(lower, h.hints_0_lower), ind_0.upper_bound(upper, h.hints_0_upper));
}
range<t_ind_0::iterator> Type::lowerUpperRange_101(const t_tuple& lower, const t_tuple& upper) const {
context h;
return lowerUpperRange_101(lower,upper,h);
}
range<t_ind_0::iterator> Type::lowerUpperRange_111(const t_tuple& lower, const t_tuple& upper, context& h) const {
t_comparator_0 comparator;
int cmp = comparator(lower, upper);
if (cmp == 0) {
    auto pos = ind_0.find(lower, h.hints_0_lower);
    auto fin = ind_0.end();
    if (pos != fin) {fin = pos; ++fin;}
    return make_range(pos, fin);
}
if (cmp > 0) {
    return make_range(ind_0.end(), ind_0.end());
}
return make_range(ind_0.lower_bound(lower, h.hints_0_lower), ind_0.upper_bound(upper, h.hints_0_upper));
}
range<t_ind_0::iterator> Type::lowerUpperRange_111(const t_tuple& lower, const t_tuple& upper) const {
context h;
return lowerUpperRange_111(lower,upper,h);
}
bool Type::empty() const {
return ind_0.empty();
}
std::vector<range<iterator>> Type::partition() const {
return ind_0.getChunks(400);
}
void Type::purge() {
ind_0.clear();
}
iterator Type::begin() const {
return ind_0.begin();
}
iterator Type::end() const {
return ind_0.end();
}
void Type::printStatistics(std::ostream& o) const {
o << " arity 3 direct b-tree index 0 lex-order [0,2,1]\n";
ind_0.printStats(o);
}
} // namespace souffle::t_btree_iii__0_2_1__101__111 
namespace souffle::t_btree_ii__0_1__11__10 {
using namespace souffle;
struct Type {
static constexpr Relation::arity_type Arity = 2;
using t_tuple = Tuple<RamDomain, 2>;
struct t_comparator_0{
 int operator()(const t_tuple& a, const t_tuple& b) const {
  return (ramBitCast<RamSigned>(a[0]) < ramBitCast<RamSigned>(b[0])) ? -1 : (ramBitCast<RamSigned>(a[0]) > ramBitCast<RamSigned>(b[0])) ? 1 :((ramBitCast<RamSigned>(a[1]) < ramBitCast<RamSigned>(b[1])) ? -1 : (ramBitCast<RamSigned>(a[1]) > ramBitCast<RamSigned>(b[1])) ? 1 :(0));
 }
bool less(const t_tuple& a, const t_tuple& b) const {
  return (ramBitCast<RamSigned>(a[0]) < ramBitCast<RamSigned>(b[0]))|| ((ramBitCast<RamSigned>(a[0]) == ramBitCast<RamSigned>(b[0])) && ((ramBitCast<RamSigned>(a[1]) < ramBitCast<RamSigned>(b[1]))));
 }
bool equal(const t_tuple& a, const t_tuple& b) const {
return (ramBitCast<RamSigned>(a[0]) == ramBitCast<RamSigned>(b[0]))&&(ramBitCast<RamSigned>(a[1]) == ramBitCast<RamSigned>(b[1]));
 }
};
using t_ind_0 = btree_set<t_tuple,t_comparator_0>;
t_ind_0 ind_0;
using iterator = t_ind_0::iterator;
struct context {
t_ind_0::operation_hints hints_0_lower;
t_ind_0::operation_hints hints_0_upper;
};
context createContext() { return context(); }
bool insert(const t_tuple& t);
bool insert(const t_tuple& t, context& h);
bool insert(const RamDomain* ramDomain);
bool insert(RamDomain a0,RamDomain a1);
bool contains(const t_tuple& t, context& h) const;
bool contains(const t_tuple& t) const;
std::size_t size() const;
iterator find(const t_tuple& t, context& h) const;
iterator find(const t_tuple& t) const;
range<iterator> lowerUpperRange_00(const t_tuple& /* lower */, const t_tuple& /* upper */, context& /* h */) const;
range<iterator> lowerUpperRange_00(const t_tuple& /* lower */, const t_tuple& /* upper */) const;
range<t_ind_0::iterator> lowerUpperRange_11(const t_tuple& lower, const t_tuple& upper, context& h) const;
range<t_ind_0::iterator> lowerUpperRange_11(const t_tuple& lower, const t_tuple& upper) const;
range<t_ind_0::iterator> lowerUpperRange_10(const t_tuple& lower, const t_tuple& upper, context& h) const;
range<t_ind_0::iterator> lowerUpperRange_10(const t_tuple& lower, const t_tuple& upper) const;
bool empty() const;
std::vector<range<iterator>> partition() const;
void purge();
iterator begin() const;
iterator end() const;
void printStatistics(std::ostream& o) const;
};
} // namespace souffle::t_btree_ii__0_1__11__10 
namespace souffle::t_btree_ii__0_1__11__10 {
using namespace souffle;
using t_ind_0 = Type::t_ind_0;
using iterator = Type::iterator;
using context = Type::context;
bool Type::insert(const t_tuple& t) {
context h;
return insert(t, h);
}
bool Type::insert(const t_tuple& t, context& h) {
if (ind_0.insert(t, h.hints_0_lower)) {
return true;
} else return false;
}
bool Type::insert(const RamDomain* ramDomain) {
RamDomain data[2];
std::copy(ramDomain, ramDomain + 2, data);
const t_tuple& tuple = reinterpret_cast<const t_tuple&>(data);
context h;
return insert(tuple, h);
}
bool Type::insert(RamDomain a0,RamDomain a1) {
RamDomain data[2] = {a0,a1};
return insert(data);
}
bool Type::contains(const t_tuple& t, context& h) const {
return ind_0.contains(t, h.hints_0_lower);
}
bool Type::contains(const t_tuple& t) const {
context h;
return contains(t, h);
}
std::size_t Type::size() const {
return ind_0.size();
}
iterator Type::find(const t_tuple& t, context& h) const {
return ind_0.find(t, h.hints_0_lower);
}
iterator Type::find(const t_tuple& t) const {
context h;
return find(t, h);
}
range<iterator> Type::lowerUpperRange_00(const t_tuple& /* lower */, const t_tuple& /* upper */, context& /* h */) const {
return range<iterator>(ind_0.begin(),ind_0.end());
}
range<iterator> Type::lowerUpperRange_00(const t_tuple& /* lower */, const t_tuple& /* upper */) const {
return range<iterator>(ind_0.begin(),ind_0.end());
}
range<t_ind_0::iterator> Type::lowerUpperRange_11(const t_tuple& lower, const t_tuple& upper, context& h) const {
t_comparator_0 comparator;
int cmp = comparator(lower, upper);
if (cmp == 0) {
    auto pos = ind_0.find(lower, h.hints_0_lower);
    auto fin = ind_0.end();
    if (pos != fin) {fin = pos; ++fin;}
    return make_range(pos, fin);
}
if (cmp > 0) {
    return make_range(ind_0.end(), ind_0.end());
}
return make_range(ind_0.lower_bound(lower, h.hints_0_lower), ind_0.upper_bound(upper, h.hints_0_upper));
}
range<t_ind_0::iterator> Type::lowerUpperRange_11(const t_tuple& lower, const t_tuple& upper) const {
context h;
return lowerUpperRange_11(lower,upper,h);
}
range<t_ind_0::iterator> Type::lowerUpperRange_10(const t_tuple& lower, const t_tuple& upper, context& h) const {
t_comparator_0 comparator;
int cmp = comparator(lower, upper);
if (cmp > 0) {
    return make_range(ind_0.end(), ind_0.end());
}
return make_range(ind_0.lower_bound(lower, h.hints_0_lower), ind_0.upper_bound(upper, h.hints_0_upper));
}
range<t_ind_0::iterator> Type::lowerUpperRange_10(const t_tuple& lower, const t_tuple& upper) const {
context h;
return lowerUpperRange_10(lower,upper,h);
}
bool Type::empty() const {
return ind_0.empty();
}
std::vector<range<iterator>> Type::partition() const {
return ind_0.getChunks(400);
}
void Type::purge() {
ind_0.clear();
}
iterator Type::begin() const {
return ind_0.begin();
}
iterator Type::end() const {
return ind_0.end();
}
void Type::printStatistics(std::ostream& o) const {
o << " arity 2 direct b-tree index 0 lex-order [0,1]\n";
ind_0.printStats(o);
}
} // namespace souffle::t_btree_ii__0_1__11__10 
namespace souffle::t_btree_ii__1_0__11__01 {
using namespace souffle;
struct Type {
static constexpr Relation::arity_type Arity = 2;
using t_tuple = Tuple<RamDomain, 2>;
struct t_comparator_0{
 int operator()(const t_tuple& a, const t_tuple& b) const {
  return (ramBitCast<RamSigned>(a[1]) < ramBitCast<RamSigned>(b[1])) ? -1 : (ramBitCast<RamSigned>(a[1]) > ramBitCast<RamSigned>(b[1])) ? 1 :((ramBitCast<RamSigned>(a[0]) < ramBitCast<RamSigned>(b[0])) ? -1 : (ramBitCast<RamSigned>(a[0]) > ramBitCast<RamSigned>(b[0])) ? 1 :(0));
 }
bool less(const t_tuple& a, const t_tuple& b) const {
  return (ramBitCast<RamSigned>(a[1]) < ramBitCast<RamSigned>(b[1]))|| ((ramBitCast<RamSigned>(a[1]) == ramBitCast<RamSigned>(b[1])) && ((ramBitCast<RamSigned>(a[0]) < ramBitCast<RamSigned>(b[0]))));
 }
bool equal(const t_tuple& a, const t_tuple& b) const {
return (ramBitCast<RamSigned>(a[1]) == ramBitCast<RamSigned>(b[1]))&&(ramBitCast<RamSigned>(a[0]) == ramBitCast<RamSigned>(b[0]));
 }
};
using t_ind_0 = btree_set<t_tuple,t_comparator_0>;
t_ind_0 ind_0;
using iterator = t_ind_0::iterator;
struct context {
t_ind_0::operation_hints hints_0_lower;
t_ind_0::operation_hints hints_0_upper;
};
context createContext() { return context(); }
bool insert(const t_tuple& t);
bool insert(const t_tuple& t, context& h);
bool insert(const RamDomain* ramDomain);
bool insert(RamDomain a0,RamDomain a1);
bool contains(const t_tuple& t, context& h) const;
bool contains(const t_tuple& t) const;
std::size_t size() const;
iterator find(const t_tuple& t, context& h) const;
iterator find(const t_tuple& t) const;
range<iterator> lowerUpperRange_00(const t_tuple& /* lower */, const t_tuple& /* upper */, context& /* h */) const;
range<iterator> lowerUpperRange_00(const t_tuple& /* lower */, const t_tuple& /* upper */) const;
range<t_ind_0::iterator> lowerUpperRange_11(const t_tuple& lower, const t_tuple& upper, context& h) const;
range<t_ind_0::iterator> lowerUpperRange_11(const t_tuple& lower, const t_tuple& upper) const;
range<t_ind_0::iterator> lowerUpperRange_01(const t_tuple& lower, const t_tuple& upper, context& h) const;
range<t_ind_0::iterator> lowerUpperRange_01(const t_tuple& lower, const t_tuple& upper) const;
bool empty() const;
std::vector<range<iterator>> partition() const;
void purge();
iterator begin() const;
iterator end() const;
void printStatistics(std::ostream& o) const;
};
} // namespace souffle::t_btree_ii__1_0__11__01 
namespace souffle::t_btree_ii__1_0__11__01 {
using namespace souffle;
using t_ind_0 = Type::t_ind_0;
using iterator = Type::iterator;
using context = Type::context;
bool Type::insert(const t_tuple& t) {
context h;
return insert(t, h);
}
bool Type::insert(const t_tuple& t, context& h) {
if (ind_0.insert(t, h.hints_0_lower)) {
return true;
} else return false;
}
bool Type::insert(const RamDomain* ramDomain) {
RamDomain data[2];
std::copy(ramDomain, ramDomain + 2, data);
const t_tuple& tuple = reinterpret_cast<const t_tuple&>(data);
context h;
return insert(tuple, h);
}
bool Type::insert(RamDomain a0,RamDomain a1) {
RamDomain data[2] = {a0,a1};
return insert(data);
}
bool Type::contains(const t_tuple& t, context& h) const {
return ind_0.contains(t, h.hints_0_lower);
}
bool Type::contains(const t_tuple& t) const {
context h;
return contains(t, h);
}
std::size_t Type::size() const {
return ind_0.size();
}
iterator Type::find(const t_tuple& t, context& h) const {
return ind_0.find(t, h.hints_0_lower);
}
iterator Type::find(const t_tuple& t) const {
context h;
return find(t, h);
}
range<iterator> Type::lowerUpperRange_00(const t_tuple& /* lower */, const t_tuple& /* upper */, context& /* h */) const {
return range<iterator>(ind_0.begin(),ind_0.end());
}
range<iterator> Type::lowerUpperRange_00(const t_tuple& /* lower */, const t_tuple& /* upper */) const {
return range<iterator>(ind_0.begin(),ind_0.end());
}
range<t_ind_0::iterator> Type::lowerUpperRange_11(const t_tuple& lower, const t_tuple& upper, context& h) const {
t_comparator_0 comparator;
int cmp = comparator(lower, upper);
if (cmp == 0) {
    auto pos = ind_0.find(lower, h.hints_0_lower);
    auto fin = ind_0.end();
    if (pos != fin) {fin = pos; ++fin;}
    return make_range(pos, fin);
}
if (cmp > 0) {
    return make_range(ind_0.end(), ind_0.end());
}
return make_range(ind_0.lower_bound(lower, h.hints_0_lower), ind_0.upper_bound(upper, h.hints_0_upper));
}
range<t_ind_0::iterator> Type::lowerUpperRange_11(const t_tuple& lower, const t_tuple& upper) const {
context h;
return lowerUpperRange_11(lower,upper,h);
}
range<t_ind_0::iterator> Type::lowerUpperRange_01(const t_tuple& lower, const t_tuple& upper, context& h) const {
t_comparator_0 comparator;
int cmp = comparator(lower, upper);
if (cmp > 0) {
    return make_range(ind_0.end(), ind_0.end());
}
return make_range(ind_0.lower_bound(lower, h.hints_0_lower), ind_0.upper_bound(upper, h.hints_0_upper));
}
range<t_ind_0::iterator> Type::lowerUpperRange_01(const t_tuple& lower, const t_tuple& upper) const {
context h;
return lowerUpperRange_01(lower,upper,h);
}
bool Type::empty() const {
return ind_0.empty();
}
std::vector<range<iterator>> Type::partition() const {
return ind_0.getChunks(400);
}
void Type::purge() {
ind_0.clear();
}
iterator Type::begin() const {
return ind_0.begin();
}
iterator Type::end() const {
return ind_0.end();
}
void Type::printStatistics(std::ostream& o) const {
o << " arity 2 direct b-tree index 0 lex-order [1,0]\n";
ind_0.printStats(o);
}
} // namespace souffle::t_btree_ii__1_0__11__01 
namespace souffle::t_btree_iui__0_1_2__110__111 {
using namespace souffle;
struct Type {
static constexpr Relation::arity_type Arity = 3;
using t_tuple = Tuple<RamDomain, 3>;
struct t_comparator_0{
 int operator()(const t_tuple& a, const t_tuple& b) const {
  return (ramBitCast<RamSigned>(a[0]) < ramBitCast<RamSigned>(b[0])) ? -1 : (ramBitCast<RamSigned>(a[0]) > ramBitCast<RamSigned>(b[0])) ? 1 :((ramBitCast<RamUnsigned>(a[1]) < ramBitCast<RamUnsigned>(b[1])) ? -1 : (ramBitCast<RamUnsigned>(a[1]) > ramBitCast<RamUnsigned>(b[1])) ? 1 :((ramBitCast<RamSigned>(a[2]) < ramBitCast<RamSigned>(b[2])) ? -1 : (ramBitCast<RamSigned>(a[2]) > ramBitCast<RamSigned>(b[2])) ? 1 :(0)));
 }
bool less(const t_tuple& a, const t_tuple& b) const {
  return (ramBitCast<RamSigned>(a[0]) < ramBitCast<RamSigned>(b[0]))|| ((ramBitCast<RamSigned>(a[0]) == ramBitCast<RamSigned>(b[0])) && ((ramBitCast<RamUnsigned>(a[1]) < ramBitCast<RamUnsigned>(b[1]))|| ((ramBitCast<RamUnsigned>(a[1]) == ramBitCast<RamUnsigned>(b[1])) && ((ramBitCast<RamSigned>(a[2]) < ramBitCast<RamSigned>(b[2]))))));
 }
bool equal(const t_tuple& a, const t_tuple& b) const {
return (ramBitCast<RamSigned>(a[0]) == ramBitCast<RamSigned>(b[0]))&&(ramBitCast<RamUnsigned>(a[1]) == ramBitCast<RamUnsigned>(b[1]))&&(ramBitCast<RamSigned>(a[2]) == ramBitCast<RamSigned>(b[2]));
 }
};
using t_ind_0 = btree_set<t_tuple,t_comparator_0>;
t_ind_0 ind_0;
using iterator = t_ind_0::iterator;
struct context {
t_ind_0::operation_hints hints_0_lower;
t_ind_0::operation_hints hints_0_upper;
};
context createContext() { return context(); }
bool insert(const t_tuple& t);
bool insert(const t_tuple& t, context& h);
bool insert(const RamDomain* ramDomain);
bool insert(RamDomain a0,RamDomain a1,RamDomain a2);
bool contains(const t_tuple& t, context& h) const;
bool contains(const t_tuple& t) const;
std::size_t size() const;
iterator find(const t_tuple& t, context& h) const;
iterator find(const t_tuple& t) const;
range<iterator> lowerUpperRange_000(const t_tuple& /* lower */, const t_tuple& /* upper */, context& /* h */) const;
range<iterator> lowerUpperRange_000(const t_tuple& /* lower */, const t_tuple& /* upper */) const;
range<t_ind_0::iterator> lowerUpperRange_110(const t_tuple& lower, const t_tuple& upper, context& h) const;
range<t_ind_0::iterator> lowerUpperRange_110(const t_tuple& lower, const t_tuple& upper) const;
range<t_ind_0::iterator> lowerUpperRange_111(const t_tuple& lower, const t_tuple& upper, context& h) const;
range<t_ind_0::iterator> lowerUpperRange_111(const t_tuple& lower, const t_tuple& upper) const;
bool empty() const;
std::vector<range<iterator>> partition() const;
void purge();
iterator begin() const;
iterator end() const;
void printStatistics(std::ostream& o) const;
};
} // namespace souffle::t_btree_iui__0_1_2__110__111 
namespace souffle::t_btree_iui__0_1_2__110__111 {
using namespace souffle;
using t_ind_0 = Type::t_ind_0;
using iterator = Type::iterator;
using context = Type::context;
bool Type::insert(const t_tuple& t) {
context h;
return insert(t, h);
}
bool Type::insert(const t_tuple& t, context& h) {
if (ind_0.insert(t, h.hints_0_lower)) {
return true;
} else return false;
}
bool Type::insert(const RamDomain* ramDomain) {
RamDomain data[3];
std::copy(ramDomain, ramDomain + 3, data);
const t_tuple& tuple = reinterpret_cast<const t_tuple&>(data);
context h;
return insert(tuple, h);
}
bool Type::insert(RamDomain a0,RamDomain a1,RamDomain a2) {
RamDomain data[3] = {a0,a1,a2};
return insert(data);
}
bool Type::contains(const t_tuple& t, context& h) const {
return ind_0.contains(t, h.hints_0_lower);
}
bool Type::contains(const t_tuple& t) const {
context h;
return contains(t, h);
}
std::size_t Type::size() const {
return ind_0.size();
}
iterator Type::find(const t_tuple& t, context& h) const {
return ind_0.find(t, h.hints_0_lower);
}
iterator Type::find(const t_tuple& t) const {
context h;
return find(t, h);
}
range<iterator> Type::lowerUpperRange_000(const t_tuple& /* lower */, const t_tuple& /* upper */, context& /* h */) const {
return range<iterator>(ind_0.begin(),ind_0.end());
}
range<iterator> Type::lowerUpperRange_000(const t_tuple& /* lower */, const t_tuple& /* upper */) const {
return range<iterator>(ind_0.begin(),ind_0.end());
}
range<t_ind_0::iterator> Type::lowerUpperRange_110(const t_tuple& lower, const t_tuple& upper, context& h) const {
t_comparator_0 comparator;
int cmp = comparator(lower, upper);
if (cmp > 0) {
    return make_range(ind_0.end(), ind_0.end());
}
return make_range(ind_0.lower_bound(lower, h.hints_0_lower), ind_0.upper_bound(upper, h.hints_0_upper));
}
range<t_ind_0::iterator> Type::lowerUpperRange_110(const t_tuple& lower, const t_tuple& upper) const {
context h;
return lowerUpperRange_110(lower,upper,h);
}
range<t_ind_0::iterator> Type::lowerUpperRange_111(const t_tuple& lower, const t_tuple& upper, context& h) const {
t_comparator_0 comparator;
int cmp = comparator(lower, upper);
if (cmp == 0) {
    auto pos = ind_0.find(lower, h.hints_0_lower);
    auto fin = ind_0.end();
    if (pos != fin) {fin = pos; ++fin;}
    return make_range(pos, fin);
}
if (cmp > 0) {
    return make_range(ind_0.end(), ind_0.end());
}
return make_range(ind_0.lower_bound(lower, h.hints_0_lower), ind_0.upper_bound(upper, h.hints_0_upper));
}
range<t_ind_0::iterator> Type::lowerUpperRange_111(const t_tuple& lower, const t_tuple& upper) const {
context h;
return lowerUpperRange_111(lower,upper,h);
}
bool Type::empty() const {
return ind_0.empty();
}
std::vector<range<iterator>> Type::partition() const {
return ind_0.getChunks(400);
}
void Type::purge() {
ind_0.clear();
}
iterator Type::begin() const {
return ind_0.begin();
}
iterator Type::end() const {
return ind_0.end();
}
void Type::printStatistics(std::ostream& o) const {
o << " arity 3 direct b-tree index 0 lex-order [0,1,2]\n";
ind_0.printStats(o);
}
} // namespace souffle::t_btree_iui__0_1_2__110__111 
namespace  souffle {
using namespace souffle;
class Stratum_AST_NodeContent_fd51b4bf60caba3f {
public:
 Stratum_AST_NodeContent_fd51b4bf60caba3f(SymbolTable& symTable,RecordTable& recordTable,ConcurrentCache<std::string,std::regex>& regexCache,bool& pruneImdtRels,bool& performIO,SignalHandler*& signalHandler,std::atomic<std::size_t>& iter,std::atomic<RamDomain>& ctr,std::string& inputDirectory,std::string& outputDirectory,t_btree_ii__0_1__11::Type& rel_AST_NodeContent_b2f3666572e60754);
void run([[maybe_unused]] const std::vector<RamDomain>& args,[[maybe_unused]] std::vector<RamDomain>& ret);
private:
SymbolTable& symTable;
RecordTable& recordTable;
ConcurrentCache<std::string,std::regex>& regexCache;
bool& pruneImdtRels;
bool& performIO;
SignalHandler*& signalHandler;
std::atomic<std::size_t>& iter;
std::atomic<RamDomain>& ctr;
std::string& inputDirectory;
std::string& outputDirectory;
t_btree_ii__0_1__11::Type* rel_AST_NodeContent_b2f3666572e60754;
};
} // namespace  souffle
namespace  souffle {
using namespace souffle;
 Stratum_AST_NodeContent_fd51b4bf60caba3f::Stratum_AST_NodeContent_fd51b4bf60caba3f(SymbolTable& symTable,RecordTable& recordTable,ConcurrentCache<std::string,std::regex>& regexCache,bool& pruneImdtRels,bool& performIO,SignalHandler*& signalHandler,std::atomic<std::size_t>& iter,std::atomic<RamDomain>& ctr,std::string& inputDirectory,std::string& outputDirectory,t_btree_ii__0_1__11::Type& rel_AST_NodeContent_b2f3666572e60754):
symTable(symTable),
recordTable(recordTable),
regexCache(regexCache),
pruneImdtRels(pruneImdtRels),
performIO(performIO),
signalHandler(signalHandler),
iter(iter),
ctr(ctr),
inputDirectory(inputDirectory),
outputDirectory(outputDirectory),
rel_AST_NodeContent_b2f3666572e60754(&rel_AST_NodeContent_b2f3666572e60754){
}

void Stratum_AST_NodeContent_fd51b4bf60caba3f::run([[maybe_unused]] const std::vector<RamDomain>& args,[[maybe_unused]] std::vector<RamDomain>& ret){
if (performIO) {
try {std::map<std::string, std::string> directiveMap({{"IO","file"},{"attributeNames","node\tcontent"},{"auxArity","0"},{"fact-dir","."},{"name","AST_NodeContent"},{"operation","input"},{"params","{\"records\": {\"AST_Location\": {\"arity\": 5, \"params\": [\"startByte\", \"startLine\", \"startColumn\", \"endLine\", \"endColumn\"]}, \"AST_Node\": {\"arity\": 2, \"params\": [\"file\", \"nodeId\"]}}, \"relation\": {\"arity\": 2, \"params\": [\"node\", \"content\"]}}"},{"types","{\"ADTs\": {}, \"records\": {\"r:AST_Location\": {\"arity\": 5, \"types\": [\"i:AST_FilePosition\", \"i:AST_FilePosition\", \"i:AST_FilePosition\", \"i:AST_FilePosition\", \"i:AST_FilePosition\"]}, \"r:AST_Node\": {\"arity\": 2, \"types\": [\"i:AST_FileId\", \"i:AST_NodeId\"]}}, \"relation\": {\"arity\": 2, \"types\": [\"r:AST_Node\", \"s:AST_Content\"]}}"}});
if (!inputDirectory.empty()) {directiveMap["fact-dir"] = inputDirectory;}
IOSystem::getInstance().getReader(directiveMap, symTable, recordTable)->readAll(*rel_AST_NodeContent_b2f3666572e60754);
} catch (std::exception& e) {std::cerr << "Error loading AST_NodeContent data: " << e.what() << '\n';
exit(1);
}
}
}

} // namespace  souffle

namespace  souffle {
using namespace souffle;
class Stratum_AST_NodeField_cc21295739297165 {
public:
 Stratum_AST_NodeField_cc21295739297165(SymbolTable& symTable,RecordTable& recordTable,ConcurrentCache<std::string,std::regex>& regexCache,bool& pruneImdtRels,bool& performIO,SignalHandler*& signalHandler,std::atomic<std::size_t>& iter,std::atomic<RamDomain>& ctr,std::string& inputDirectory,std::string& outputDirectory,t_btree_iii__0_2_1__101__111::Type& rel_AST_NodeField_ca02670731ce3c99);
void run([[maybe_unused]] const std::vector<RamDomain>& args,[[maybe_unused]] std::vector<RamDomain>& ret);
private:
SymbolTable& symTable;
RecordTable& recordTable;
ConcurrentCache<std::string,std::regex>& regexCache;
bool& pruneImdtRels;
bool& performIO;
SignalHandler*& signalHandler;
std::atomic<std::size_t>& iter;
std::atomic<RamDomain>& ctr;
std::string& inputDirectory;
std::string& outputDirectory;
t_btree_iii__0_2_1__101__111::Type* rel_AST_NodeField_ca02670731ce3c99;
};
} // namespace  souffle
namespace  souffle {
using namespace souffle;
 Stratum_AST_NodeField_cc21295739297165::Stratum_AST_NodeField_cc21295739297165(SymbolTable& symTable,RecordTable& recordTable,ConcurrentCache<std::string,std::regex>& regexCache,bool& pruneImdtRels,bool& performIO,SignalHandler*& signalHandler,std::atomic<std::size_t>& iter,std::atomic<RamDomain>& ctr,std::string& inputDirectory,std::string& outputDirectory,t_btree_iii__0_2_1__101__111::Type& rel_AST_NodeField_ca02670731ce3c99):
symTable(symTable),
recordTable(recordTable),
regexCache(regexCache),
pruneImdtRels(pruneImdtRels),
performIO(performIO),
signalHandler(signalHandler),
iter(iter),
ctr(ctr),
inputDirectory(inputDirectory),
outputDirectory(outputDirectory),
rel_AST_NodeField_ca02670731ce3c99(&rel_AST_NodeField_ca02670731ce3c99){
}

void Stratum_AST_NodeField_cc21295739297165::run([[maybe_unused]] const std::vector<RamDomain>& args,[[maybe_unused]] std::vector<RamDomain>& ret){
if (performIO) {
try {std::map<std::string, std::string> directiveMap({{"IO","file"},{"attributeNames","parent\tchild\tfield"},{"auxArity","0"},{"fact-dir","."},{"name","AST_NodeField"},{"operation","input"},{"params","{\"records\": {\"AST_Location\": {\"arity\": 5, \"params\": [\"startByte\", \"startLine\", \"startColumn\", \"endLine\", \"endColumn\"]}, \"AST_Node\": {\"arity\": 2, \"params\": [\"file\", \"nodeId\"]}}, \"relation\": {\"arity\": 3, \"params\": [\"parent\", \"child\", \"field\"]}}"},{"types","{\"ADTs\": {}, \"records\": {\"r:AST_Location\": {\"arity\": 5, \"types\": [\"i:AST_FilePosition\", \"i:AST_FilePosition\", \"i:AST_FilePosition\", \"i:AST_FilePosition\", \"i:AST_FilePosition\"]}, \"r:AST_Node\": {\"arity\": 2, \"types\": [\"i:AST_FileId\", \"i:AST_NodeId\"]}}, \"relation\": {\"arity\": 3, \"types\": [\"r:AST_Node\", \"r:AST_Node\", \"s:AST_Field\"]}}"}});
if (!inputDirectory.empty()) {directiveMap["fact-dir"] = inputDirectory;}
IOSystem::getInstance().getReader(directiveMap, symTable, recordTable)->readAll(*rel_AST_NodeField_ca02670731ce3c99);
} catch (std::exception& e) {std::cerr << "Error loading AST_NodeField data: " << e.what() << '\n';
exit(1);
}
}
}

} // namespace  souffle

namespace  souffle {
using namespace souffle;
class Stratum_AST_NodeLocation_89d765aa14237a09 {
public:
 Stratum_AST_NodeLocation_89d765aa14237a09(SymbolTable& symTable,RecordTable& recordTable,ConcurrentCache<std::string,std::regex>& regexCache,bool& pruneImdtRels,bool& performIO,SignalHandler*& signalHandler,std::atomic<std::size_t>& iter,std::atomic<RamDomain>& ctr,std::string& inputDirectory,std::string& outputDirectory,t_btree_ii__0_1__11__10::Type& rel_AST_NodeLocation_5f3f38ee7a82c12a);
void run([[maybe_unused]] const std::vector<RamDomain>& args,[[maybe_unused]] std::vector<RamDomain>& ret);
private:
SymbolTable& symTable;
RecordTable& recordTable;
ConcurrentCache<std::string,std::regex>& regexCache;
bool& pruneImdtRels;
bool& performIO;
SignalHandler*& signalHandler;
std::atomic<std::size_t>& iter;
std::atomic<RamDomain>& ctr;
std::string& inputDirectory;
std::string& outputDirectory;
t_btree_ii__0_1__11__10::Type* rel_AST_NodeLocation_5f3f38ee7a82c12a;
};
} // namespace  souffle
namespace  souffle {
using namespace souffle;
 Stratum_AST_NodeLocation_89d765aa14237a09::Stratum_AST_NodeLocation_89d765aa14237a09(SymbolTable& symTable,RecordTable& recordTable,ConcurrentCache<std::string,std::regex>& regexCache,bool& pruneImdtRels,bool& performIO,SignalHandler*& signalHandler,std::atomic<std::size_t>& iter,std::atomic<RamDomain>& ctr,std::string& inputDirectory,std::string& outputDirectory,t_btree_ii__0_1__11__10::Type& rel_AST_NodeLocation_5f3f38ee7a82c12a):
symTable(symTable),
recordTable(recordTable),
regexCache(regexCache),
pruneImdtRels(pruneImdtRels),
performIO(performIO),
signalHandler(signalHandler),
iter(iter),
ctr(ctr),
inputDirectory(inputDirectory),
outputDirectory(outputDirectory),
rel_AST_NodeLocation_5f3f38ee7a82c12a(&rel_AST_NodeLocation_5f3f38ee7a82c12a){
}

void Stratum_AST_NodeLocation_89d765aa14237a09::run([[maybe_unused]] const std::vector<RamDomain>& args,[[maybe_unused]] std::vector<RamDomain>& ret){
if (performIO) {
try {std::map<std::string, std::string> directiveMap({{"IO","file"},{"attributeNames","node\tlocation"},{"auxArity","0"},{"fact-dir","."},{"name","AST_NodeLocation"},{"operation","input"},{"params","{\"records\": {\"AST_Location\": {\"arity\": 5, \"params\": [\"startByte\", \"startLine\", \"startColumn\", \"endLine\", \"endColumn\"]}, \"AST_Node\": {\"arity\": 2, \"params\": [\"file\", \"nodeId\"]}}, \"relation\": {\"arity\": 2, \"params\": [\"node\", \"location\"]}}"},{"types","{\"ADTs\": {}, \"records\": {\"r:AST_Location\": {\"arity\": 5, \"types\": [\"i:AST_FilePosition\", \"i:AST_FilePosition\", \"i:AST_FilePosition\", \"i:AST_FilePosition\", \"i:AST_FilePosition\"]}, \"r:AST_Node\": {\"arity\": 2, \"types\": [\"i:AST_FileId\", \"i:AST_NodeId\"]}}, \"relation\": {\"arity\": 2, \"types\": [\"r:AST_Node\", \"r:AST_Location\"]}}"}});
if (!inputDirectory.empty()) {directiveMap["fact-dir"] = inputDirectory;}
IOSystem::getInstance().getReader(directiveMap, symTable, recordTable)->readAll(*rel_AST_NodeLocation_5f3f38ee7a82c12a);
} catch (std::exception& e) {std::cerr << "Error loading AST_NodeLocation data: " << e.what() << '\n';
exit(1);
}
}
}

} // namespace  souffle

namespace  souffle {
using namespace souffle;
class Stratum_AST_NodeType_400775685fb3e630 {
public:
 Stratum_AST_NodeType_400775685fb3e630(SymbolTable& symTable,RecordTable& recordTable,ConcurrentCache<std::string,std::regex>& regexCache,bool& pruneImdtRels,bool& performIO,SignalHandler*& signalHandler,std::atomic<std::size_t>& iter,std::atomic<RamDomain>& ctr,std::string& inputDirectory,std::string& outputDirectory,t_btree_ii__1_0__11__01::Type& rel_AST_NodeType_b38285ae9991409e);
void run([[maybe_unused]] const std::vector<RamDomain>& args,[[maybe_unused]] std::vector<RamDomain>& ret);
private:
SymbolTable& symTable;
RecordTable& recordTable;
ConcurrentCache<std::string,std::regex>& regexCache;
bool& pruneImdtRels;
bool& performIO;
SignalHandler*& signalHandler;
std::atomic<std::size_t>& iter;
std::atomic<RamDomain>& ctr;
std::string& inputDirectory;
std::string& outputDirectory;
t_btree_ii__1_0__11__01::Type* rel_AST_NodeType_b38285ae9991409e;
};
} // namespace  souffle
namespace  souffle {
using namespace souffle;
 Stratum_AST_NodeType_400775685fb3e630::Stratum_AST_NodeType_400775685fb3e630(SymbolTable& symTable,RecordTable& recordTable,ConcurrentCache<std::string,std::regex>& regexCache,bool& pruneImdtRels,bool& performIO,SignalHandler*& signalHandler,std::atomic<std::size_t>& iter,std::atomic<RamDomain>& ctr,std::string& inputDirectory,std::string& outputDirectory,t_btree_ii__1_0__11__01::Type& rel_AST_NodeType_b38285ae9991409e):
symTable(symTable),
recordTable(recordTable),
regexCache(regexCache),
pruneImdtRels(pruneImdtRels),
performIO(performIO),
signalHandler(signalHandler),
iter(iter),
ctr(ctr),
inputDirectory(inputDirectory),
outputDirectory(outputDirectory),
rel_AST_NodeType_b38285ae9991409e(&rel_AST_NodeType_b38285ae9991409e){
}

void Stratum_AST_NodeType_400775685fb3e630::run([[maybe_unused]] const std::vector<RamDomain>& args,[[maybe_unused]] std::vector<RamDomain>& ret){
if (performIO) {
try {std::map<std::string, std::string> directiveMap({{"IO","file"},{"attributeNames","node\ttype"},{"auxArity","0"},{"fact-dir","."},{"name","AST_NodeType"},{"operation","input"},{"params","{\"records\": {\"AST_Location\": {\"arity\": 5, \"params\": [\"startByte\", \"startLine\", \"startColumn\", \"endLine\", \"endColumn\"]}, \"AST_Node\": {\"arity\": 2, \"params\": [\"file\", \"nodeId\"]}}, \"relation\": {\"arity\": 2, \"params\": [\"node\", \"type\"]}}"},{"types","{\"ADTs\": {}, \"records\": {\"r:AST_Location\": {\"arity\": 5, \"types\": [\"i:AST_FilePosition\", \"i:AST_FilePosition\", \"i:AST_FilePosition\", \"i:AST_FilePosition\", \"i:AST_FilePosition\"]}, \"r:AST_Node\": {\"arity\": 2, \"types\": [\"i:AST_FileId\", \"i:AST_NodeId\"]}}, \"relation\": {\"arity\": 2, \"types\": [\"r:AST_Node\", \"s:AST_Type\"]}}"}});
if (!inputDirectory.empty()) {directiveMap["fact-dir"] = inputDirectory;}
IOSystem::getInstance().getReader(directiveMap, symTable, recordTable)->readAll(*rel_AST_NodeType_b38285ae9991409e);
} catch (std::exception& e) {std::cerr << "Error loading AST_NodeType data: " << e.what() << '\n';
exit(1);
}
}
}

} // namespace  souffle

namespace  souffle {
using namespace souffle;
class Stratum_AST_ParentChild_798cb83c96de8e4d {
public:
 Stratum_AST_ParentChild_798cb83c96de8e4d(SymbolTable& symTable,RecordTable& recordTable,ConcurrentCache<std::string,std::regex>& regexCache,bool& pruneImdtRels,bool& performIO,SignalHandler*& signalHandler,std::atomic<std::size_t>& iter,std::atomic<RamDomain>& ctr,std::string& inputDirectory,std::string& outputDirectory,t_btree_iui__0_1_2__110__111::Type& rel_AST_ParentChild_be6259205eb66578);
void run([[maybe_unused]] const std::vector<RamDomain>& args,[[maybe_unused]] std::vector<RamDomain>& ret);
private:
SymbolTable& symTable;
RecordTable& recordTable;
ConcurrentCache<std::string,std::regex>& regexCache;
bool& pruneImdtRels;
bool& performIO;
SignalHandler*& signalHandler;
std::atomic<std::size_t>& iter;
std::atomic<RamDomain>& ctr;
std::string& inputDirectory;
std::string& outputDirectory;
t_btree_iui__0_1_2__110__111::Type* rel_AST_ParentChild_be6259205eb66578;
};
} // namespace  souffle
namespace  souffle {
using namespace souffle;
 Stratum_AST_ParentChild_798cb83c96de8e4d::Stratum_AST_ParentChild_798cb83c96de8e4d(SymbolTable& symTable,RecordTable& recordTable,ConcurrentCache<std::string,std::regex>& regexCache,bool& pruneImdtRels,bool& performIO,SignalHandler*& signalHandler,std::atomic<std::size_t>& iter,std::atomic<RamDomain>& ctr,std::string& inputDirectory,std::string& outputDirectory,t_btree_iui__0_1_2__110__111::Type& rel_AST_ParentChild_be6259205eb66578):
symTable(symTable),
recordTable(recordTable),
regexCache(regexCache),
pruneImdtRels(pruneImdtRels),
performIO(performIO),
signalHandler(signalHandler),
iter(iter),
ctr(ctr),
inputDirectory(inputDirectory),
outputDirectory(outputDirectory),
rel_AST_ParentChild_be6259205eb66578(&rel_AST_ParentChild_be6259205eb66578){
}

void Stratum_AST_ParentChild_798cb83c96de8e4d::run([[maybe_unused]] const std::vector<RamDomain>& args,[[maybe_unused]] std::vector<RamDomain>& ret){
if (performIO) {
try {std::map<std::string, std::string> directiveMap({{"IO","file"},{"attributeNames","parent\tindex\tchild"},{"auxArity","0"},{"fact-dir","."},{"name","AST_ParentChild"},{"operation","input"},{"params","{\"records\": {\"AST_Location\": {\"arity\": 5, \"params\": [\"startByte\", \"startLine\", \"startColumn\", \"endLine\", \"endColumn\"]}, \"AST_Node\": {\"arity\": 2, \"params\": [\"file\", \"nodeId\"]}}, \"relation\": {\"arity\": 3, \"params\": [\"parent\", \"index\", \"child\"]}}"},{"types","{\"ADTs\": {}, \"records\": {\"r:AST_Location\": {\"arity\": 5, \"types\": [\"i:AST_FilePosition\", \"i:AST_FilePosition\", \"i:AST_FilePosition\", \"i:AST_FilePosition\", \"i:AST_FilePosition\"]}, \"r:AST_Node\": {\"arity\": 2, \"types\": [\"i:AST_FileId\", \"i:AST_NodeId\"]}}, \"relation\": {\"arity\": 3, \"types\": [\"r:AST_Node\", \"u:unsigned\", \"r:AST_Node\"]}}"}});
if (!inputDirectory.empty()) {directiveMap["fact-dir"] = inputDirectory;}
IOSystem::getInstance().getReader(directiveMap, symTable, recordTable)->readAll(*rel_AST_ParentChild_be6259205eb66578);
} catch (std::exception& e) {std::cerr << "Error loading AST_ParentChild data: " << e.what() << '\n';
exit(1);
}
}
}

} // namespace  souffle

namespace  souffle {
using namespace souffle;
class Stratum_Rule_0e6b7aa9ece342e5 {
public:
 Stratum_Rule_0e6b7aa9ece342e5(SymbolTable& symTable,RecordTable& recordTable,ConcurrentCache<std::string,std::regex>& regexCache,bool& pruneImdtRels,bool& performIO,SignalHandler*& signalHandler,std::atomic<std::size_t>& iter,std::atomic<RamDomain>& ctr,std::string& inputDirectory,std::string& outputDirectory,t_btree_ii__0_1__11::Type& rel_AST_NodeContent_b2f3666572e60754,t_btree_iii__0_2_1__101__111::Type& rel_AST_NodeField_ca02670731ce3c99,t_btree_ii__1_0__11__01::Type& rel_AST_NodeType_b38285ae9991409e,t_btree_iui__0_1_2__110__111::Type& rel_AST_ParentChild_be6259205eb66578,t_btree_ii__0_1__11::Type& rel_Rule_159a590c7b2c8d95);
void run([[maybe_unused]] const std::vector<RamDomain>& args,[[maybe_unused]] std::vector<RamDomain>& ret);
private:
SymbolTable& symTable;
RecordTable& recordTable;
ConcurrentCache<std::string,std::regex>& regexCache;
bool& pruneImdtRels;
bool& performIO;
SignalHandler*& signalHandler;
std::atomic<std::size_t>& iter;
std::atomic<RamDomain>& ctr;
std::string& inputDirectory;
std::string& outputDirectory;
t_btree_ii__0_1__11::Type* rel_AST_NodeContent_b2f3666572e60754;
t_btree_iii__0_2_1__101__111::Type* rel_AST_NodeField_ca02670731ce3c99;
t_btree_ii__1_0__11__01::Type* rel_AST_NodeType_b38285ae9991409e;
t_btree_iui__0_1_2__110__111::Type* rel_AST_ParentChild_be6259205eb66578;
t_btree_ii__0_1__11::Type* rel_Rule_159a590c7b2c8d95;
};
} // namespace  souffle
namespace  souffle {
using namespace souffle;
 Stratum_Rule_0e6b7aa9ece342e5::Stratum_Rule_0e6b7aa9ece342e5(SymbolTable& symTable,RecordTable& recordTable,ConcurrentCache<std::string,std::regex>& regexCache,bool& pruneImdtRels,bool& performIO,SignalHandler*& signalHandler,std::atomic<std::size_t>& iter,std::atomic<RamDomain>& ctr,std::string& inputDirectory,std::string& outputDirectory,t_btree_ii__0_1__11::Type& rel_AST_NodeContent_b2f3666572e60754,t_btree_iii__0_2_1__101__111::Type& rel_AST_NodeField_ca02670731ce3c99,t_btree_ii__1_0__11__01::Type& rel_AST_NodeType_b38285ae9991409e,t_btree_iui__0_1_2__110__111::Type& rel_AST_ParentChild_be6259205eb66578,t_btree_ii__0_1__11::Type& rel_Rule_159a590c7b2c8d95):
symTable(symTable),
recordTable(recordTable),
regexCache(regexCache),
pruneImdtRels(pruneImdtRels),
performIO(performIO),
signalHandler(signalHandler),
iter(iter),
ctr(ctr),
inputDirectory(inputDirectory),
outputDirectory(outputDirectory),
rel_AST_NodeContent_b2f3666572e60754(&rel_AST_NodeContent_b2f3666572e60754),
rel_AST_NodeField_ca02670731ce3c99(&rel_AST_NodeField_ca02670731ce3c99),
rel_AST_NodeType_b38285ae9991409e(&rel_AST_NodeType_b38285ae9991409e),
rel_AST_ParentChild_be6259205eb66578(&rel_AST_ParentChild_be6259205eb66578),
rel_Rule_159a590c7b2c8d95(&rel_Rule_159a590c7b2c8d95){
}

void Stratum_Rule_0e6b7aa9ece342e5::run([[maybe_unused]] const std::vector<RamDomain>& args,[[maybe_unused]] std::vector<RamDomain>& ret){
signalHandler->setMsg(R"_(Rule("sql_lang_create_table",node1) :- 
   AST_NodeType(node1,"call"),
   !AST_NodeField(node1,_,"receiver"),
   !AST_NodeField(node1,_,"block"),
   AST_NodeField(node1,node2,"method"),
   AST_NodeType(node2,"constant"),
   AST_NodeContent(node2,"CREATE"),
   AST_NodeField(node1,node3,"arguments"),
   AST_NodeType(node3,"argument_list"),
   AST_ParentChild(node3,0,node4),
   AST_NodeType(node4,"call"),
   !AST_NodeField(node4,_,"receiver"),
   !AST_NodeField(node4,_,"block"),
   AST_NodeField(node4,node5,"method"),
   AST_NodeType(node5,"constant"),
   AST_NodeContent(node5,"TABLE"),
   AST_NodeField(node4,node6,"arguments"),
   AST_NodeType(node6,"argument_list"),
   AST_ParentChild(node6,0,node7),
   AST_NodeType(node7,"call"),
   !AST_NodeField(node7,_,"block"),
   AST_NodeField(node7,node8,"receiver"),
   AST_NodeType(node8,"identifier"),
   AST_NodeContent(node8,"public"),
   AST_ParentChild(node7,0,node9),
   AST_NodeType(node9,"."),
   AST_ParentChild(node9,0,node10),
   AST_NodeType(node10,"identifier"),
   AST_NodeContent(node10,""),
   AST_NodeField(node7,node11,"arguments"),
   AST_NodeType(node11,"argument_list"),
   AST_ParentChild(node11,0,node12),
   AST_NodeType(node12,"call"),
   !AST_NodeField(node12,_,"receiver"),
   !AST_NodeField(node12,_,"block"),
   AST_NodeField(node12,node13,"method"),
   AST_NodeType(node13,"global_variable"),
   AST_NodeField(node12,node14,"arguments"),
   AST_NodeType(node14,"argument_list"),
   AST_ParentChild(node14,0,node15),
   AST_NodeType(node15,"parenthesized_statements"),
   AST_ParentChild(node15,0,node16),
   AST_NodeType(node16,"("),
   AST_ParentChild(node15,1,node17),
   AST_NodeType(node17,"ERROR"),
   AST_ParentChild(node17,0,node18),
   AST_NodeType(node18,"<"),
   AST_ParentChild(node15,2,node19),
   AST_NodeType(node19,"global_variable"),
   AST_ParentChild(node15,3,node20),
   AST_NodeType(node20,"ERROR"),
   AST_ParentChild(node20,0,node21),
   AST_NodeType(node21,">"),
   AST_ParentChild(node15,4,node22),
   AST_NodeType(node22,")").
in file rules.dl [1:1-1:1954])_");
if(!(rel_AST_ParentChild_be6259205eb66578->empty()) && !(rel_AST_NodeField_ca02670731ce3c99->empty()) && !(rel_AST_NodeType_b38285ae9991409e->empty()) && !(rel_AST_NodeContent_b2f3666572e60754->empty())) {
[&](){
CREATE_OP_CONTEXT(rel_AST_NodeContent_b2f3666572e60754_op_ctxt,rel_AST_NodeContent_b2f3666572e60754->createContext());
CREATE_OP_CONTEXT(rel_AST_NodeField_ca02670731ce3c99_op_ctxt,rel_AST_NodeField_ca02670731ce3c99->createContext());
CREATE_OP_CONTEXT(rel_AST_NodeType_b38285ae9991409e_op_ctxt,rel_AST_NodeType_b38285ae9991409e->createContext());
CREATE_OP_CONTEXT(rel_AST_ParentChild_be6259205eb66578_op_ctxt,rel_AST_ParentChild_be6259205eb66578->createContext());
CREATE_OP_CONTEXT(rel_Rule_159a590c7b2c8d95_op_ctxt,rel_Rule_159a590c7b2c8d95->createContext());
auto range = rel_AST_NodeType_b38285ae9991409e->lowerUpperRange_01(Tuple<RamDomain,2>{{ramBitCast<RamDomain>(MIN_RAM_SIGNED), ramBitCast(RamSigned(0))}},Tuple<RamDomain,2>{{ramBitCast<RamDomain>(MAX_RAM_SIGNED), ramBitCast(RamSigned(0))}},READ_OP_CONTEXT(rel_AST_NodeType_b38285ae9991409e_op_ctxt));
for(const auto& env0 : range) {
if( !(!rel_AST_NodeField_ca02670731ce3c99->lowerUpperRange_101(Tuple<RamDomain,3>{{ramBitCast(env0[0]), ramBitCast<RamDomain>(MIN_RAM_SIGNED), ramBitCast(RamSigned(1))}},Tuple<RamDomain,3>{{ramBitCast(env0[0]), ramBitCast<RamDomain>(MAX_RAM_SIGNED), ramBitCast(RamSigned(1))}},READ_OP_CONTEXT(rel_AST_NodeField_ca02670731ce3c99_op_ctxt)).empty()) && !(!rel_AST_NodeField_ca02670731ce3c99->lowerUpperRange_101(Tuple<RamDomain,3>{{ramBitCast(env0[0]), ramBitCast<RamDomain>(MIN_RAM_SIGNED), ramBitCast(RamSigned(2))}},Tuple<RamDomain,3>{{ramBitCast(env0[0]), ramBitCast<RamDomain>(MAX_RAM_SIGNED), ramBitCast(RamSigned(2))}},READ_OP_CONTEXT(rel_AST_NodeField_ca02670731ce3c99_op_ctxt)).empty())) {
auto range = rel_AST_NodeField_ca02670731ce3c99->lowerUpperRange_101(Tuple<RamDomain,3>{{ramBitCast(env0[0]), ramBitCast<RamDomain>(MIN_RAM_SIGNED), ramBitCast(RamSigned(3))}},Tuple<RamDomain,3>{{ramBitCast(env0[0]), ramBitCast<RamDomain>(MAX_RAM_SIGNED), ramBitCast(RamSigned(3))}},READ_OP_CONTEXT(rel_AST_NodeField_ca02670731ce3c99_op_ctxt));
for(const auto& env1 : range) {
if( rel_AST_NodeContent_b2f3666572e60754->contains(Tuple<RamDomain,2>{{ramBitCast(env1[1]),ramBitCast(RamSigned(4))}},READ_OP_CONTEXT(rel_AST_NodeContent_b2f3666572e60754_op_ctxt)) && rel_AST_NodeType_b38285ae9991409e->contains(Tuple<RamDomain,2>{{ramBitCast(env1[1]),ramBitCast(RamSigned(5))}},READ_OP_CONTEXT(rel_AST_NodeType_b38285ae9991409e_op_ctxt))) {
auto range = rel_AST_NodeField_ca02670731ce3c99->lowerUpperRange_101(Tuple<RamDomain,3>{{ramBitCast(env0[0]), ramBitCast<RamDomain>(MIN_RAM_SIGNED), ramBitCast(RamSigned(6))}},Tuple<RamDomain,3>{{ramBitCast(env0[0]), ramBitCast<RamDomain>(MAX_RAM_SIGNED), ramBitCast(RamSigned(6))}},READ_OP_CONTEXT(rel_AST_NodeField_ca02670731ce3c99_op_ctxt));
for(const auto& env2 : range) {
if( rel_AST_NodeType_b38285ae9991409e->contains(Tuple<RamDomain,2>{{ramBitCast(env2[1]),ramBitCast(RamSigned(7))}},READ_OP_CONTEXT(rel_AST_NodeType_b38285ae9991409e_op_ctxt))) {
auto range = rel_AST_ParentChild_be6259205eb66578->lowerUpperRange_110(Tuple<RamDomain,3>{{ramBitCast(env2[1]), ramBitCast(RamUnsigned(0)), ramBitCast<RamDomain>(MIN_RAM_SIGNED)}},Tuple<RamDomain,3>{{ramBitCast(env2[1]), ramBitCast(RamUnsigned(0)), ramBitCast<RamDomain>(MAX_RAM_SIGNED)}},READ_OP_CONTEXT(rel_AST_ParentChild_be6259205eb66578_op_ctxt));
for(const auto& env3 : range) {
if( rel_AST_NodeType_b38285ae9991409e->contains(Tuple<RamDomain,2>{{ramBitCast(env3[2]),ramBitCast(RamSigned(0))}},READ_OP_CONTEXT(rel_AST_NodeType_b38285ae9991409e_op_ctxt)) && !(!rel_AST_NodeField_ca02670731ce3c99->lowerUpperRange_101(Tuple<RamDomain,3>{{ramBitCast(env3[2]), ramBitCast<RamDomain>(MIN_RAM_SIGNED), ramBitCast(RamSigned(1))}},Tuple<RamDomain,3>{{ramBitCast(env3[2]), ramBitCast<RamDomain>(MAX_RAM_SIGNED), ramBitCast(RamSigned(1))}},READ_OP_CONTEXT(rel_AST_NodeField_ca02670731ce3c99_op_ctxt)).empty()) && !(!rel_AST_NodeField_ca02670731ce3c99->lowerUpperRange_101(Tuple<RamDomain,3>{{ramBitCast(env3[2]), ramBitCast<RamDomain>(MIN_RAM_SIGNED), ramBitCast(RamSigned(2))}},Tuple<RamDomain,3>{{ramBitCast(env3[2]), ramBitCast<RamDomain>(MAX_RAM_SIGNED), ramBitCast(RamSigned(2))}},READ_OP_CONTEXT(rel_AST_NodeField_ca02670731ce3c99_op_ctxt)).empty())) {
auto range = rel_AST_NodeField_ca02670731ce3c99->lowerUpperRange_101(Tuple<RamDomain,3>{{ramBitCast(env3[2]), ramBitCast<RamDomain>(MIN_RAM_SIGNED), ramBitCast(RamSigned(3))}},Tuple<RamDomain,3>{{ramBitCast(env3[2]), ramBitCast<RamDomain>(MAX_RAM_SIGNED), ramBitCast(RamSigned(3))}},READ_OP_CONTEXT(rel_AST_NodeField_ca02670731ce3c99_op_ctxt));
for(const auto& env4 : range) {
if( rel_AST_NodeContent_b2f3666572e60754->contains(Tuple<RamDomain,2>{{ramBitCast(env4[1]),ramBitCast(RamSigned(8))}},READ_OP_CONTEXT(rel_AST_NodeContent_b2f3666572e60754_op_ctxt)) && rel_AST_NodeType_b38285ae9991409e->contains(Tuple<RamDomain,2>{{ramBitCast(env4[1]),ramBitCast(RamSigned(5))}},READ_OP_CONTEXT(rel_AST_NodeType_b38285ae9991409e_op_ctxt))) {
auto range = rel_AST_NodeField_ca02670731ce3c99->lowerUpperRange_101(Tuple<RamDomain,3>{{ramBitCast(env3[2]), ramBitCast<RamDomain>(MIN_RAM_SIGNED), ramBitCast(RamSigned(6))}},Tuple<RamDomain,3>{{ramBitCast(env3[2]), ramBitCast<RamDomain>(MAX_RAM_SIGNED), ramBitCast(RamSigned(6))}},READ_OP_CONTEXT(rel_AST_NodeField_ca02670731ce3c99_op_ctxt));
for(const auto& env5 : range) {
if( rel_AST_NodeType_b38285ae9991409e->contains(Tuple<RamDomain,2>{{ramBitCast(env5[1]),ramBitCast(RamSigned(7))}},READ_OP_CONTEXT(rel_AST_NodeType_b38285ae9991409e_op_ctxt))) {
auto range = rel_AST_ParentChild_be6259205eb66578->lowerUpperRange_110(Tuple<RamDomain,3>{{ramBitCast(env5[1]), ramBitCast(RamUnsigned(0)), ramBitCast<RamDomain>(MIN_RAM_SIGNED)}},Tuple<RamDomain,3>{{ramBitCast(env5[1]), ramBitCast(RamUnsigned(0)), ramBitCast<RamDomain>(MAX_RAM_SIGNED)}},READ_OP_CONTEXT(rel_AST_ParentChild_be6259205eb66578_op_ctxt));
for(const auto& env6 : range) {
if( rel_AST_NodeType_b38285ae9991409e->contains(Tuple<RamDomain,2>{{ramBitCast(env6[2]),ramBitCast(RamSigned(0))}},READ_OP_CONTEXT(rel_AST_NodeType_b38285ae9991409e_op_ctxt)) && !(!rel_AST_NodeField_ca02670731ce3c99->lowerUpperRange_101(Tuple<RamDomain,3>{{ramBitCast(env6[2]), ramBitCast<RamDomain>(MIN_RAM_SIGNED), ramBitCast(RamSigned(2))}},Tuple<RamDomain,3>{{ramBitCast(env6[2]), ramBitCast<RamDomain>(MAX_RAM_SIGNED), ramBitCast(RamSigned(2))}},READ_OP_CONTEXT(rel_AST_NodeField_ca02670731ce3c99_op_ctxt)).empty())) {
auto range = rel_AST_NodeField_ca02670731ce3c99->lowerUpperRange_101(Tuple<RamDomain,3>{{ramBitCast(env6[2]), ramBitCast<RamDomain>(MIN_RAM_SIGNED), ramBitCast(RamSigned(1))}},Tuple<RamDomain,3>{{ramBitCast(env6[2]), ramBitCast<RamDomain>(MAX_RAM_SIGNED), ramBitCast(RamSigned(1))}},READ_OP_CONTEXT(rel_AST_NodeField_ca02670731ce3c99_op_ctxt));
for(const auto& env7 : range) {
if( rel_AST_NodeContent_b2f3666572e60754->contains(Tuple<RamDomain,2>{{ramBitCast(env7[1]),ramBitCast(RamSigned(9))}},READ_OP_CONTEXT(rel_AST_NodeContent_b2f3666572e60754_op_ctxt)) && rel_AST_NodeType_b38285ae9991409e->contains(Tuple<RamDomain,2>{{ramBitCast(env7[1]),ramBitCast(RamSigned(10))}},READ_OP_CONTEXT(rel_AST_NodeType_b38285ae9991409e_op_ctxt))) {
auto range = rel_AST_ParentChild_be6259205eb66578->lowerUpperRange_110(Tuple<RamDomain,3>{{ramBitCast(env6[2]), ramBitCast(RamUnsigned(0)), ramBitCast<RamDomain>(MIN_RAM_SIGNED)}},Tuple<RamDomain,3>{{ramBitCast(env6[2]), ramBitCast(RamUnsigned(0)), ramBitCast<RamDomain>(MAX_RAM_SIGNED)}},READ_OP_CONTEXT(rel_AST_ParentChild_be6259205eb66578_op_ctxt));
for(const auto& env8 : range) {
if( rel_AST_NodeType_b38285ae9991409e->contains(Tuple<RamDomain,2>{{ramBitCast(env8[2]),ramBitCast(RamSigned(11))}},READ_OP_CONTEXT(rel_AST_NodeType_b38285ae9991409e_op_ctxt))) {
auto range = rel_AST_ParentChild_be6259205eb66578->lowerUpperRange_110(Tuple<RamDomain,3>{{ramBitCast(env8[2]), ramBitCast(RamUnsigned(0)), ramBitCast<RamDomain>(MIN_RAM_SIGNED)}},Tuple<RamDomain,3>{{ramBitCast(env8[2]), ramBitCast(RamUnsigned(0)), ramBitCast<RamDomain>(MAX_RAM_SIGNED)}},READ_OP_CONTEXT(rel_AST_ParentChild_be6259205eb66578_op_ctxt));
for(const auto& env9 : range) {
if( rel_AST_NodeContent_b2f3666572e60754->contains(Tuple<RamDomain,2>{{ramBitCast(env9[2]),ramBitCast(RamSigned(12))}},READ_OP_CONTEXT(rel_AST_NodeContent_b2f3666572e60754_op_ctxt)) && rel_AST_NodeType_b38285ae9991409e->contains(Tuple<RamDomain,2>{{ramBitCast(env9[2]),ramBitCast(RamSigned(10))}},READ_OP_CONTEXT(rel_AST_NodeType_b38285ae9991409e_op_ctxt))) {
auto range = rel_AST_NodeField_ca02670731ce3c99->lowerUpperRange_101(Tuple<RamDomain,3>{{ramBitCast(env6[2]), ramBitCast<RamDomain>(MIN_RAM_SIGNED), ramBitCast(RamSigned(6))}},Tuple<RamDomain,3>{{ramBitCast(env6[2]), ramBitCast<RamDomain>(MAX_RAM_SIGNED), ramBitCast(RamSigned(6))}},READ_OP_CONTEXT(rel_AST_NodeField_ca02670731ce3c99_op_ctxt));
for(const auto& env10 : range) {
if( rel_AST_NodeType_b38285ae9991409e->contains(Tuple<RamDomain,2>{{ramBitCast(env10[1]),ramBitCast(RamSigned(7))}},READ_OP_CONTEXT(rel_AST_NodeType_b38285ae9991409e_op_ctxt))) {
auto range = rel_AST_ParentChild_be6259205eb66578->lowerUpperRange_110(Tuple<RamDomain,3>{{ramBitCast(env10[1]), ramBitCast(RamUnsigned(0)), ramBitCast<RamDomain>(MIN_RAM_SIGNED)}},Tuple<RamDomain,3>{{ramBitCast(env10[1]), ramBitCast(RamUnsigned(0)), ramBitCast<RamDomain>(MAX_RAM_SIGNED)}},READ_OP_CONTEXT(rel_AST_ParentChild_be6259205eb66578_op_ctxt));
for(const auto& env11 : range) {
if( rel_AST_NodeType_b38285ae9991409e->contains(Tuple<RamDomain,2>{{ramBitCast(env11[2]),ramBitCast(RamSigned(0))}},READ_OP_CONTEXT(rel_AST_NodeType_b38285ae9991409e_op_ctxt)) && !(!rel_AST_NodeField_ca02670731ce3c99->lowerUpperRange_101(Tuple<RamDomain,3>{{ramBitCast(env11[2]), ramBitCast<RamDomain>(MIN_RAM_SIGNED), ramBitCast(RamSigned(1))}},Tuple<RamDomain,3>{{ramBitCast(env11[2]), ramBitCast<RamDomain>(MAX_RAM_SIGNED), ramBitCast(RamSigned(1))}},READ_OP_CONTEXT(rel_AST_NodeField_ca02670731ce3c99_op_ctxt)).empty()) && !(!rel_AST_NodeField_ca02670731ce3c99->lowerUpperRange_101(Tuple<RamDomain,3>{{ramBitCast(env11[2]), ramBitCast<RamDomain>(MIN_RAM_SIGNED), ramBitCast(RamSigned(2))}},Tuple<RamDomain,3>{{ramBitCast(env11[2]), ramBitCast<RamDomain>(MAX_RAM_SIGNED), ramBitCast(RamSigned(2))}},READ_OP_CONTEXT(rel_AST_NodeField_ca02670731ce3c99_op_ctxt)).empty())) {
auto range = rel_AST_NodeField_ca02670731ce3c99->lowerUpperRange_101(Tuple<RamDomain,3>{{ramBitCast(env11[2]), ramBitCast<RamDomain>(MIN_RAM_SIGNED), ramBitCast(RamSigned(3))}},Tuple<RamDomain,3>{{ramBitCast(env11[2]), ramBitCast<RamDomain>(MAX_RAM_SIGNED), ramBitCast(RamSigned(3))}},READ_OP_CONTEXT(rel_AST_NodeField_ca02670731ce3c99_op_ctxt));
for(const auto& env12 : range) {
if( rel_AST_NodeType_b38285ae9991409e->contains(Tuple<RamDomain,2>{{ramBitCast(env12[1]),ramBitCast(RamSigned(13))}},READ_OP_CONTEXT(rel_AST_NodeType_b38285ae9991409e_op_ctxt))) {
auto range = rel_AST_NodeField_ca02670731ce3c99->lowerUpperRange_101(Tuple<RamDomain,3>{{ramBitCast(env11[2]), ramBitCast<RamDomain>(MIN_RAM_SIGNED), ramBitCast(RamSigned(6))}},Tuple<RamDomain,3>{{ramBitCast(env11[2]), ramBitCast<RamDomain>(MAX_RAM_SIGNED), ramBitCast(RamSigned(6))}},READ_OP_CONTEXT(rel_AST_NodeField_ca02670731ce3c99_op_ctxt));
for(const auto& env13 : range) {
if( rel_AST_NodeType_b38285ae9991409e->contains(Tuple<RamDomain,2>{{ramBitCast(env13[1]),ramBitCast(RamSigned(7))}},READ_OP_CONTEXT(rel_AST_NodeType_b38285ae9991409e_op_ctxt))) {
auto range = rel_AST_ParentChild_be6259205eb66578->lowerUpperRange_110(Tuple<RamDomain,3>{{ramBitCast(env13[1]), ramBitCast(RamUnsigned(0)), ramBitCast<RamDomain>(MIN_RAM_SIGNED)}},Tuple<RamDomain,3>{{ramBitCast(env13[1]), ramBitCast(RamUnsigned(0)), ramBitCast<RamDomain>(MAX_RAM_SIGNED)}},READ_OP_CONTEXT(rel_AST_ParentChild_be6259205eb66578_op_ctxt));
for(const auto& env14 : range) {
if( rel_AST_NodeType_b38285ae9991409e->contains(Tuple<RamDomain,2>{{ramBitCast(env14[2]),ramBitCast(RamSigned(14))}},READ_OP_CONTEXT(rel_AST_NodeType_b38285ae9991409e_op_ctxt))) {
auto range = rel_AST_ParentChild_be6259205eb66578->lowerUpperRange_110(Tuple<RamDomain,3>{{ramBitCast(env14[2]), ramBitCast(RamUnsigned(0)), ramBitCast<RamDomain>(MIN_RAM_SIGNED)}},Tuple<RamDomain,3>{{ramBitCast(env14[2]), ramBitCast(RamUnsigned(0)), ramBitCast<RamDomain>(MAX_RAM_SIGNED)}},READ_OP_CONTEXT(rel_AST_ParentChild_be6259205eb66578_op_ctxt));
for(const auto& env15 : range) {
if( rel_AST_NodeType_b38285ae9991409e->contains(Tuple<RamDomain,2>{{ramBitCast(env15[2]),ramBitCast(RamSigned(15))}},READ_OP_CONTEXT(rel_AST_NodeType_b38285ae9991409e_op_ctxt))) {
auto range = rel_AST_ParentChild_be6259205eb66578->lowerUpperRange_110(Tuple<RamDomain,3>{{ramBitCast(env14[2]), ramBitCast(RamUnsigned(1)), ramBitCast<RamDomain>(MIN_RAM_SIGNED)}},Tuple<RamDomain,3>{{ramBitCast(env14[2]), ramBitCast(RamUnsigned(1)), ramBitCast<RamDomain>(MAX_RAM_SIGNED)}},READ_OP_CONTEXT(rel_AST_ParentChild_be6259205eb66578_op_ctxt));
for(const auto& env16 : range) {
if( rel_AST_NodeType_b38285ae9991409e->contains(Tuple<RamDomain,2>{{ramBitCast(env16[2]),ramBitCast(RamSigned(16))}},READ_OP_CONTEXT(rel_AST_NodeType_b38285ae9991409e_op_ctxt))) {
auto range = rel_AST_ParentChild_be6259205eb66578->lowerUpperRange_110(Tuple<RamDomain,3>{{ramBitCast(env16[2]), ramBitCast(RamUnsigned(0)), ramBitCast<RamDomain>(MIN_RAM_SIGNED)}},Tuple<RamDomain,3>{{ramBitCast(env16[2]), ramBitCast(RamUnsigned(0)), ramBitCast<RamDomain>(MAX_RAM_SIGNED)}},READ_OP_CONTEXT(rel_AST_ParentChild_be6259205eb66578_op_ctxt));
for(const auto& env17 : range) {
if( rel_AST_NodeType_b38285ae9991409e->contains(Tuple<RamDomain,2>{{ramBitCast(env17[2]),ramBitCast(RamSigned(17))}},READ_OP_CONTEXT(rel_AST_NodeType_b38285ae9991409e_op_ctxt))) {
auto range = rel_AST_ParentChild_be6259205eb66578->lowerUpperRange_110(Tuple<RamDomain,3>{{ramBitCast(env14[2]), ramBitCast(RamUnsigned(2)), ramBitCast<RamDomain>(MIN_RAM_SIGNED)}},Tuple<RamDomain,3>{{ramBitCast(env14[2]), ramBitCast(RamUnsigned(2)), ramBitCast<RamDomain>(MAX_RAM_SIGNED)}},READ_OP_CONTEXT(rel_AST_ParentChild_be6259205eb66578_op_ctxt));
for(const auto& env18 : range) {
if( rel_AST_NodeType_b38285ae9991409e->contains(Tuple<RamDomain,2>{{ramBitCast(env18[2]),ramBitCast(RamSigned(13))}},READ_OP_CONTEXT(rel_AST_NodeType_b38285ae9991409e_op_ctxt))) {
auto range = rel_AST_ParentChild_be6259205eb66578->lowerUpperRange_110(Tuple<RamDomain,3>{{ramBitCast(env14[2]), ramBitCast(RamUnsigned(3)), ramBitCast<RamDomain>(MIN_RAM_SIGNED)}},Tuple<RamDomain,3>{{ramBitCast(env14[2]), ramBitCast(RamUnsigned(3)), ramBitCast<RamDomain>(MAX_RAM_SIGNED)}},READ_OP_CONTEXT(rel_AST_ParentChild_be6259205eb66578_op_ctxt));
for(const auto& env19 : range) {
if( rel_AST_NodeType_b38285ae9991409e->contains(Tuple<RamDomain,2>{{ramBitCast(env19[2]),ramBitCast(RamSigned(16))}},READ_OP_CONTEXT(rel_AST_NodeType_b38285ae9991409e_op_ctxt))) {
auto range = rel_AST_ParentChild_be6259205eb66578->lowerUpperRange_110(Tuple<RamDomain,3>{{ramBitCast(env19[2]), ramBitCast(RamUnsigned(0)), ramBitCast<RamDomain>(MIN_RAM_SIGNED)}},Tuple<RamDomain,3>{{ramBitCast(env19[2]), ramBitCast(RamUnsigned(0)), ramBitCast<RamDomain>(MAX_RAM_SIGNED)}},READ_OP_CONTEXT(rel_AST_ParentChild_be6259205eb66578_op_ctxt));
for(const auto& env20 : range) {
if( rel_AST_NodeType_b38285ae9991409e->contains(Tuple<RamDomain,2>{{ramBitCast(env20[2]),ramBitCast(RamSigned(18))}},READ_OP_CONTEXT(rel_AST_NodeType_b38285ae9991409e_op_ctxt))) {
auto range = rel_AST_ParentChild_be6259205eb66578->lowerUpperRange_110(Tuple<RamDomain,3>{{ramBitCast(env14[2]), ramBitCast(RamUnsigned(4)), ramBitCast<RamDomain>(MIN_RAM_SIGNED)}},Tuple<RamDomain,3>{{ramBitCast(env14[2]), ramBitCast(RamUnsigned(4)), ramBitCast<RamDomain>(MAX_RAM_SIGNED)}},READ_OP_CONTEXT(rel_AST_ParentChild_be6259205eb66578_op_ctxt));
for(const auto& env21 : range) {
if( rel_AST_NodeType_b38285ae9991409e->contains(Tuple<RamDomain,2>{{ramBitCast(env21[2]),ramBitCast(RamSigned(19))}},READ_OP_CONTEXT(rel_AST_NodeType_b38285ae9991409e_op_ctxt))) {
Tuple<RamDomain,2> tuple{{ramBitCast(RamSigned(20)),ramBitCast(env0[0])}};
rel_Rule_159a590c7b2c8d95->insert(tuple,READ_OP_CONTEXT(rel_Rule_159a590c7b2c8d95_op_ctxt));
break;
}
}
break;
}
}
}
}
break;
}
}
break;
}
}
}
}
break;
}
}
}
}
}
}
break;
}
}
}
}
}
}
}
}
}
}
}
}
}
}
}
}
}
}
}
}
}
}
}
}
}
}
}
();}
signalHandler->setMsg(R"_(Rule("sql_lang_create_table",node1) :- 
   AST_NodeType(node1,"call"),
   !AST_NodeField(node1,_,"receiver"),
   !AST_NodeField(node1,_,"block"),
   AST_NodeField(node1,node2,"method"),
   AST_NodeType(node2,"constant"),
   AST_NodeContent(node2,"CREATE"),
   AST_NodeField(node1,node3,"arguments"),
   AST_NodeType(node3,"argument_list"),
   AST_ParentChild(node3,0,node4),
   AST_NodeType(node4,"call"),
   !AST_NodeField(node4,_,"receiver"),
   !AST_NodeField(node4,_,"block"),
   AST_NodeField(node4,node5,"method"),
   AST_NodeType(node5,"constant"),
   AST_NodeContent(node5,"TABLE"),
   AST_NodeField(node4,node6,"arguments"),
   AST_NodeType(node6,"argument_list"),
   AST_ParentChild(node6,0,node7),
   AST_NodeType(node7,"call"),
   !AST_NodeField(node7,_,"receiver"),
   !AST_NodeField(node7,_,"block"),
   AST_NodeField(node7,node8,"method"),
   AST_NodeType(node8,"global_variable"),
   AST_NodeField(node7,node9,"arguments"),
   AST_NodeType(node9,"argument_list"),
   AST_ParentChild(node9,0,node10),
   AST_NodeType(node10,"parenthesized_statements"),
   AST_ParentChild(node10,0,node11),
   AST_NodeType(node11,"("),
   AST_ParentChild(node10,1,node12),
   AST_NodeType(node12,"ERROR"),
   AST_ParentChild(node12,0,node13),
   AST_NodeType(node13,"<"),
   AST_ParentChild(node10,2,node14),
   AST_NodeType(node14,"global_variable"),
   AST_ParentChild(node10,3,node15),
   AST_NodeType(node15,"ERROR"),
   AST_ParentChild(node15,0,node16),
   AST_NodeType(node16,">"),
   AST_ParentChild(node10,4,node17),
   AST_NodeType(node17,")").
in file rules.dl [2:1-2:1492])_");
if(!(rel_AST_ParentChild_be6259205eb66578->empty()) && !(rel_AST_NodeField_ca02670731ce3c99->empty()) && !(rel_AST_NodeType_b38285ae9991409e->empty()) && !(rel_AST_NodeContent_b2f3666572e60754->empty())) {
[&](){
CREATE_OP_CONTEXT(rel_AST_NodeContent_b2f3666572e60754_op_ctxt,rel_AST_NodeContent_b2f3666572e60754->createContext());
CREATE_OP_CONTEXT(rel_AST_NodeField_ca02670731ce3c99_op_ctxt,rel_AST_NodeField_ca02670731ce3c99->createContext());
CREATE_OP_CONTEXT(rel_AST_NodeType_b38285ae9991409e_op_ctxt,rel_AST_NodeType_b38285ae9991409e->createContext());
CREATE_OP_CONTEXT(rel_AST_ParentChild_be6259205eb66578_op_ctxt,rel_AST_ParentChild_be6259205eb66578->createContext());
CREATE_OP_CONTEXT(rel_Rule_159a590c7b2c8d95_op_ctxt,rel_Rule_159a590c7b2c8d95->createContext());
auto range = rel_AST_NodeType_b38285ae9991409e->lowerUpperRange_01(Tuple<RamDomain,2>{{ramBitCast<RamDomain>(MIN_RAM_SIGNED), ramBitCast(RamSigned(0))}},Tuple<RamDomain,2>{{ramBitCast<RamDomain>(MAX_RAM_SIGNED), ramBitCast(RamSigned(0))}},READ_OP_CONTEXT(rel_AST_NodeType_b38285ae9991409e_op_ctxt));
for(const auto& env0 : range) {
if( !(!rel_AST_NodeField_ca02670731ce3c99->lowerUpperRange_101(Tuple<RamDomain,3>{{ramBitCast(env0[0]), ramBitCast<RamDomain>(MIN_RAM_SIGNED), ramBitCast(RamSigned(1))}},Tuple<RamDomain,3>{{ramBitCast(env0[0]), ramBitCast<RamDomain>(MAX_RAM_SIGNED), ramBitCast(RamSigned(1))}},READ_OP_CONTEXT(rel_AST_NodeField_ca02670731ce3c99_op_ctxt)).empty()) && !(!rel_AST_NodeField_ca02670731ce3c99->lowerUpperRange_101(Tuple<RamDomain,3>{{ramBitCast(env0[0]), ramBitCast<RamDomain>(MIN_RAM_SIGNED), ramBitCast(RamSigned(2))}},Tuple<RamDomain,3>{{ramBitCast(env0[0]), ramBitCast<RamDomain>(MAX_RAM_SIGNED), ramBitCast(RamSigned(2))}},READ_OP_CONTEXT(rel_AST_NodeField_ca02670731ce3c99_op_ctxt)).empty())) {
auto range = rel_AST_NodeField_ca02670731ce3c99->lowerUpperRange_101(Tuple<RamDomain,3>{{ramBitCast(env0[0]), ramBitCast<RamDomain>(MIN_RAM_SIGNED), ramBitCast(RamSigned(3))}},Tuple<RamDomain,3>{{ramBitCast(env0[0]), ramBitCast<RamDomain>(MAX_RAM_SIGNED), ramBitCast(RamSigned(3))}},READ_OP_CONTEXT(rel_AST_NodeField_ca02670731ce3c99_op_ctxt));
for(const auto& env1 : range) {
if( rel_AST_NodeContent_b2f3666572e60754->contains(Tuple<RamDomain,2>{{ramBitCast(env1[1]),ramBitCast(RamSigned(4))}},READ_OP_CONTEXT(rel_AST_NodeContent_b2f3666572e60754_op_ctxt)) && rel_AST_NodeType_b38285ae9991409e->contains(Tuple<RamDomain,2>{{ramBitCast(env1[1]),ramBitCast(RamSigned(5))}},READ_OP_CONTEXT(rel_AST_NodeType_b38285ae9991409e_op_ctxt))) {
auto range = rel_AST_NodeField_ca02670731ce3c99->lowerUpperRange_101(Tuple<RamDomain,3>{{ramBitCast(env0[0]), ramBitCast<RamDomain>(MIN_RAM_SIGNED), ramBitCast(RamSigned(6))}},Tuple<RamDomain,3>{{ramBitCast(env0[0]), ramBitCast<RamDomain>(MAX_RAM_SIGNED), ramBitCast(RamSigned(6))}},READ_OP_CONTEXT(rel_AST_NodeField_ca02670731ce3c99_op_ctxt));
for(const auto& env2 : range) {
if( rel_AST_NodeType_b38285ae9991409e->contains(Tuple<RamDomain,2>{{ramBitCast(env2[1]),ramBitCast(RamSigned(7))}},READ_OP_CONTEXT(rel_AST_NodeType_b38285ae9991409e_op_ctxt))) {
auto range = rel_AST_ParentChild_be6259205eb66578->lowerUpperRange_110(Tuple<RamDomain,3>{{ramBitCast(env2[1]), ramBitCast(RamUnsigned(0)), ramBitCast<RamDomain>(MIN_RAM_SIGNED)}},Tuple<RamDomain,3>{{ramBitCast(env2[1]), ramBitCast(RamUnsigned(0)), ramBitCast<RamDomain>(MAX_RAM_SIGNED)}},READ_OP_CONTEXT(rel_AST_ParentChild_be6259205eb66578_op_ctxt));
for(const auto& env3 : range) {
if( rel_AST_NodeType_b38285ae9991409e->contains(Tuple<RamDomain,2>{{ramBitCast(env3[2]),ramBitCast(RamSigned(0))}},READ_OP_CONTEXT(rel_AST_NodeType_b38285ae9991409e_op_ctxt)) && !(!rel_AST_NodeField_ca02670731ce3c99->lowerUpperRange_101(Tuple<RamDomain,3>{{ramBitCast(env3[2]), ramBitCast<RamDomain>(MIN_RAM_SIGNED), ramBitCast(RamSigned(1))}},Tuple<RamDomain,3>{{ramBitCast(env3[2]), ramBitCast<RamDomain>(MAX_RAM_SIGNED), ramBitCast(RamSigned(1))}},READ_OP_CONTEXT(rel_AST_NodeField_ca02670731ce3c99_op_ctxt)).empty()) && !(!rel_AST_NodeField_ca02670731ce3c99->lowerUpperRange_101(Tuple<RamDomain,3>{{ramBitCast(env3[2]), ramBitCast<RamDomain>(MIN_RAM_SIGNED), ramBitCast(RamSigned(2))}},Tuple<RamDomain,3>{{ramBitCast(env3[2]), ramBitCast<RamDomain>(MAX_RAM_SIGNED), ramBitCast(RamSigned(2))}},READ_OP_CONTEXT(rel_AST_NodeField_ca02670731ce3c99_op_ctxt)).empty())) {
auto range = rel_AST_NodeField_ca02670731ce3c99->lowerUpperRange_101(Tuple<RamDomain,3>{{ramBitCast(env3[2]), ramBitCast<RamDomain>(MIN_RAM_SIGNED), ramBitCast(RamSigned(3))}},Tuple<RamDomain,3>{{ramBitCast(env3[2]), ramBitCast<RamDomain>(MAX_RAM_SIGNED), ramBitCast(RamSigned(3))}},READ_OP_CONTEXT(rel_AST_NodeField_ca02670731ce3c99_op_ctxt));
for(const auto& env4 : range) {
if( rel_AST_NodeContent_b2f3666572e60754->contains(Tuple<RamDomain,2>{{ramBitCast(env4[1]),ramBitCast(RamSigned(8))}},READ_OP_CONTEXT(rel_AST_NodeContent_b2f3666572e60754_op_ctxt)) && rel_AST_NodeType_b38285ae9991409e->contains(Tuple<RamDomain,2>{{ramBitCast(env4[1]),ramBitCast(RamSigned(5))}},READ_OP_CONTEXT(rel_AST_NodeType_b38285ae9991409e_op_ctxt))) {
auto range = rel_AST_NodeField_ca02670731ce3c99->lowerUpperRange_101(Tuple<RamDomain,3>{{ramBitCast(env3[2]), ramBitCast<RamDomain>(MIN_RAM_SIGNED), ramBitCast(RamSigned(6))}},Tuple<RamDomain,3>{{ramBitCast(env3[2]), ramBitCast<RamDomain>(MAX_RAM_SIGNED), ramBitCast(RamSigned(6))}},READ_OP_CONTEXT(rel_AST_NodeField_ca02670731ce3c99_op_ctxt));
for(const auto& env5 : range) {
if( rel_AST_NodeType_b38285ae9991409e->contains(Tuple<RamDomain,2>{{ramBitCast(env5[1]),ramBitCast(RamSigned(7))}},READ_OP_CONTEXT(rel_AST_NodeType_b38285ae9991409e_op_ctxt))) {
auto range = rel_AST_ParentChild_be6259205eb66578->lowerUpperRange_110(Tuple<RamDomain,3>{{ramBitCast(env5[1]), ramBitCast(RamUnsigned(0)), ramBitCast<RamDomain>(MIN_RAM_SIGNED)}},Tuple<RamDomain,3>{{ramBitCast(env5[1]), ramBitCast(RamUnsigned(0)), ramBitCast<RamDomain>(MAX_RAM_SIGNED)}},READ_OP_CONTEXT(rel_AST_ParentChild_be6259205eb66578_op_ctxt));
for(const auto& env6 : range) {
if( rel_AST_NodeType_b38285ae9991409e->contains(Tuple<RamDomain,2>{{ramBitCast(env6[2]),ramBitCast(RamSigned(0))}},READ_OP_CONTEXT(rel_AST_NodeType_b38285ae9991409e_op_ctxt)) && !(!rel_AST_NodeField_ca02670731ce3c99->lowerUpperRange_101(Tuple<RamDomain,3>{{ramBitCast(env6[2]), ramBitCast<RamDomain>(MIN_RAM_SIGNED), ramBitCast(RamSigned(1))}},Tuple<RamDomain,3>{{ramBitCast(env6[2]), ramBitCast<RamDomain>(MAX_RAM_SIGNED), ramBitCast(RamSigned(1))}},READ_OP_CONTEXT(rel_AST_NodeField_ca02670731ce3c99_op_ctxt)).empty()) && !(!rel_AST_NodeField_ca02670731ce3c99->lowerUpperRange_101(Tuple<RamDomain,3>{{ramBitCast(env6[2]), ramBitCast<RamDomain>(MIN_RAM_SIGNED), ramBitCast(RamSigned(2))}},Tuple<RamDomain,3>{{ramBitCast(env6[2]), ramBitCast<RamDomain>(MAX_RAM_SIGNED), ramBitCast(RamSigned(2))}},READ_OP_CONTEXT(rel_AST_NodeField_ca02670731ce3c99_op_ctxt)).empty())) {
auto range = rel_AST_NodeField_ca02670731ce3c99->lowerUpperRange_101(Tuple<RamDomain,3>{{ramBitCast(env6[2]), ramBitCast<RamDomain>(MIN_RAM_SIGNED), ramBitCast(RamSigned(3))}},Tuple<RamDomain,3>{{ramBitCast(env6[2]), ramBitCast<RamDomain>(MAX_RAM_SIGNED), ramBitCast(RamSigned(3))}},READ_OP_CONTEXT(rel_AST_NodeField_ca02670731ce3c99_op_ctxt));
for(const auto& env7 : range) {
if( rel_AST_NodeType_b38285ae9991409e->contains(Tuple<RamDomain,2>{{ramBitCast(env7[1]),ramBitCast(RamSigned(13))}},READ_OP_CONTEXT(rel_AST_NodeType_b38285ae9991409e_op_ctxt))) {
auto range = rel_AST_NodeField_ca02670731ce3c99->lowerUpperRange_101(Tuple<RamDomain,3>{{ramBitCast(env6[2]), ramBitCast<RamDomain>(MIN_RAM_SIGNED), ramBitCast(RamSigned(6))}},Tuple<RamDomain,3>{{ramBitCast(env6[2]), ramBitCast<RamDomain>(MAX_RAM_SIGNED), ramBitCast(RamSigned(6))}},READ_OP_CONTEXT(rel_AST_NodeField_ca02670731ce3c99_op_ctxt));
for(const auto& env8 : range) {
if( rel_AST_NodeType_b38285ae9991409e->contains(Tuple<RamDomain,2>{{ramBitCast(env8[1]),ramBitCast(RamSigned(7))}},READ_OP_CONTEXT(rel_AST_NodeType_b38285ae9991409e_op_ctxt))) {
auto range = rel_AST_ParentChild_be6259205eb66578->lowerUpperRange_110(Tuple<RamDomain,3>{{ramBitCast(env8[1]), ramBitCast(RamUnsigned(0)), ramBitCast<RamDomain>(MIN_RAM_SIGNED)}},Tuple<RamDomain,3>{{ramBitCast(env8[1]), ramBitCast(RamUnsigned(0)), ramBitCast<RamDomain>(MAX_RAM_SIGNED)}},READ_OP_CONTEXT(rel_AST_ParentChild_be6259205eb66578_op_ctxt));
for(const auto& env9 : range) {
if( rel_AST_NodeType_b38285ae9991409e->contains(Tuple<RamDomain,2>{{ramBitCast(env9[2]),ramBitCast(RamSigned(14))}},READ_OP_CONTEXT(rel_AST_NodeType_b38285ae9991409e_op_ctxt))) {
auto range = rel_AST_ParentChild_be6259205eb66578->lowerUpperRange_110(Tuple<RamDomain,3>{{ramBitCast(env9[2]), ramBitCast(RamUnsigned(0)), ramBitCast<RamDomain>(MIN_RAM_SIGNED)}},Tuple<RamDomain,3>{{ramBitCast(env9[2]), ramBitCast(RamUnsigned(0)), ramBitCast<RamDomain>(MAX_RAM_SIGNED)}},READ_OP_CONTEXT(rel_AST_ParentChild_be6259205eb66578_op_ctxt));
for(const auto& env10 : range) {
if( rel_AST_NodeType_b38285ae9991409e->contains(Tuple<RamDomain,2>{{ramBitCast(env10[2]),ramBitCast(RamSigned(15))}},READ_OP_CONTEXT(rel_AST_NodeType_b38285ae9991409e_op_ctxt))) {
auto range = rel_AST_ParentChild_be6259205eb66578->lowerUpperRange_110(Tuple<RamDomain,3>{{ramBitCast(env9[2]), ramBitCast(RamUnsigned(1)), ramBitCast<RamDomain>(MIN_RAM_SIGNED)}},Tuple<RamDomain,3>{{ramBitCast(env9[2]), ramBitCast(RamUnsigned(1)), ramBitCast<RamDomain>(MAX_RAM_SIGNED)}},READ_OP_CONTEXT(rel_AST_ParentChild_be6259205eb66578_op_ctxt));
for(const auto& env11 : range) {
if( rel_AST_NodeType_b38285ae9991409e->contains(Tuple<RamDomain,2>{{ramBitCast(env11[2]),ramBitCast(RamSigned(16))}},READ_OP_CONTEXT(rel_AST_NodeType_b38285ae9991409e_op_ctxt))) {
auto range = rel_AST_ParentChild_be6259205eb66578->lowerUpperRange_110(Tuple<RamDomain,3>{{ramBitCast(env11[2]), ramBitCast(RamUnsigned(0)), ramBitCast<RamDomain>(MIN_RAM_SIGNED)}},Tuple<RamDomain,3>{{ramBitCast(env11[2]), ramBitCast(RamUnsigned(0)), ramBitCast<RamDomain>(MAX_RAM_SIGNED)}},READ_OP_CONTEXT(rel_AST_ParentChild_be6259205eb66578_op_ctxt));
for(const auto& env12 : range) {
if( rel_AST_NodeType_b38285ae9991409e->contains(Tuple<RamDomain,2>{{ramBitCast(env12[2]),ramBitCast(RamSigned(17))}},READ_OP_CONTEXT(rel_AST_NodeType_b38285ae9991409e_op_ctxt))) {
auto range = rel_AST_ParentChild_be6259205eb66578->lowerUpperRange_110(Tuple<RamDomain,3>{{ramBitCast(env9[2]), ramBitCast(RamUnsigned(2)), ramBitCast<RamDomain>(MIN_RAM_SIGNED)}},Tuple<RamDomain,3>{{ramBitCast(env9[2]), ramBitCast(RamUnsigned(2)), ramBitCast<RamDomain>(MAX_RAM_SIGNED)}},READ_OP_CONTEXT(rel_AST_ParentChild_be6259205eb66578_op_ctxt));
for(const auto& env13 : range) {
if( rel_AST_NodeType_b38285ae9991409e->contains(Tuple<RamDomain,2>{{ramBitCast(env13[2]),ramBitCast(RamSigned(13))}},READ_OP_CONTEXT(rel_AST_NodeType_b38285ae9991409e_op_ctxt))) {
auto range = rel_AST_ParentChild_be6259205eb66578->lowerUpperRange_110(Tuple<RamDomain,3>{{ramBitCast(env9[2]), ramBitCast(RamUnsigned(3)), ramBitCast<RamDomain>(MIN_RAM_SIGNED)}},Tuple<RamDomain,3>{{ramBitCast(env9[2]), ramBitCast(RamUnsigned(3)), ramBitCast<RamDomain>(MAX_RAM_SIGNED)}},READ_OP_CONTEXT(rel_AST_ParentChild_be6259205eb66578_op_ctxt));
for(const auto& env14 : range) {
if( rel_AST_NodeType_b38285ae9991409e->contains(Tuple<RamDomain,2>{{ramBitCast(env14[2]),ramBitCast(RamSigned(16))}},READ_OP_CONTEXT(rel_AST_NodeType_b38285ae9991409e_op_ctxt))) {
auto range = rel_AST_ParentChild_be6259205eb66578->lowerUpperRange_110(Tuple<RamDomain,3>{{ramBitCast(env14[2]), ramBitCast(RamUnsigned(0)), ramBitCast<RamDomain>(MIN_RAM_SIGNED)}},Tuple<RamDomain,3>{{ramBitCast(env14[2]), ramBitCast(RamUnsigned(0)), ramBitCast<RamDomain>(MAX_RAM_SIGNED)}},READ_OP_CONTEXT(rel_AST_ParentChild_be6259205eb66578_op_ctxt));
for(const auto& env15 : range) {
if( rel_AST_NodeType_b38285ae9991409e->contains(Tuple<RamDomain,2>{{ramBitCast(env15[2]),ramBitCast(RamSigned(18))}},READ_OP_CONTEXT(rel_AST_NodeType_b38285ae9991409e_op_ctxt))) {
auto range = rel_AST_ParentChild_be6259205eb66578->lowerUpperRange_110(Tuple<RamDomain,3>{{ramBitCast(env9[2]), ramBitCast(RamUnsigned(4)), ramBitCast<RamDomain>(MIN_RAM_SIGNED)}},Tuple<RamDomain,3>{{ramBitCast(env9[2]), ramBitCast(RamUnsigned(4)), ramBitCast<RamDomain>(MAX_RAM_SIGNED)}},READ_OP_CONTEXT(rel_AST_ParentChild_be6259205eb66578_op_ctxt));
for(const auto& env16 : range) {
if( rel_AST_NodeType_b38285ae9991409e->contains(Tuple<RamDomain,2>{{ramBitCast(env16[2]),ramBitCast(RamSigned(19))}},READ_OP_CONTEXT(rel_AST_NodeType_b38285ae9991409e_op_ctxt))) {
Tuple<RamDomain,2> tuple{{ramBitCast(RamSigned(20)),ramBitCast(env0[0])}};
rel_Rule_159a590c7b2c8d95->insert(tuple,READ_OP_CONTEXT(rel_Rule_159a590c7b2c8d95_op_ctxt));
break;
}
}
break;
}
}
}
}
break;
}
}
break;
}
}
}
}
break;
}
}
}
}
}
}
break;
}
}
}
}
}
}
}
}
}
}
}
}
}
}
}
}
}
();}
if (pruneImdtRels) rel_AST_NodeContent_b2f3666572e60754->purge();
if (pruneImdtRels) rel_AST_NodeField_ca02670731ce3c99->purge();
if (pruneImdtRels) rel_AST_NodeType_b38285ae9991409e->purge();
if (pruneImdtRels) rel_AST_ParentChild_be6259205eb66578->purge();
}

} // namespace  souffle

namespace  souffle {
using namespace souffle;
class Stratum_RuleMatch_4394a245605301ce {
public:
 Stratum_RuleMatch_4394a245605301ce(SymbolTable& symTable,RecordTable& recordTable,ConcurrentCache<std::string,std::regex>& regexCache,bool& pruneImdtRels,bool& performIO,SignalHandler*& signalHandler,std::atomic<std::size_t>& iter,std::atomic<RamDomain>& ctr,std::string& inputDirectory,std::string& outputDirectory,t_btree_ii__0_1__11__10::Type& rel_AST_NodeLocation_5f3f38ee7a82c12a,t_btree_ii__0_1__11::Type& rel_Rule_159a590c7b2c8d95,t_btree_ii__0_1__11::Type& rel_RuleMatch_8974a5cadf2d4779);
void run([[maybe_unused]] const std::vector<RamDomain>& args,[[maybe_unused]] std::vector<RamDomain>& ret);
private:
SymbolTable& symTable;
RecordTable& recordTable;
ConcurrentCache<std::string,std::regex>& regexCache;
bool& pruneImdtRels;
bool& performIO;
SignalHandler*& signalHandler;
std::atomic<std::size_t>& iter;
std::atomic<RamDomain>& ctr;
std::string& inputDirectory;
std::string& outputDirectory;
t_btree_ii__0_1__11__10::Type* rel_AST_NodeLocation_5f3f38ee7a82c12a;
t_btree_ii__0_1__11::Type* rel_Rule_159a590c7b2c8d95;
t_btree_ii__0_1__11::Type* rel_RuleMatch_8974a5cadf2d4779;
};
} // namespace  souffle
namespace  souffle {
using namespace souffle;
 Stratum_RuleMatch_4394a245605301ce::Stratum_RuleMatch_4394a245605301ce(SymbolTable& symTable,RecordTable& recordTable,ConcurrentCache<std::string,std::regex>& regexCache,bool& pruneImdtRels,bool& performIO,SignalHandler*& signalHandler,std::atomic<std::size_t>& iter,std::atomic<RamDomain>& ctr,std::string& inputDirectory,std::string& outputDirectory,t_btree_ii__0_1__11__10::Type& rel_AST_NodeLocation_5f3f38ee7a82c12a,t_btree_ii__0_1__11::Type& rel_Rule_159a590c7b2c8d95,t_btree_ii__0_1__11::Type& rel_RuleMatch_8974a5cadf2d4779):
symTable(symTable),
recordTable(recordTable),
regexCache(regexCache),
pruneImdtRels(pruneImdtRels),
performIO(performIO),
signalHandler(signalHandler),
iter(iter),
ctr(ctr),
inputDirectory(inputDirectory),
outputDirectory(outputDirectory),
rel_AST_NodeLocation_5f3f38ee7a82c12a(&rel_AST_NodeLocation_5f3f38ee7a82c12a),
rel_Rule_159a590c7b2c8d95(&rel_Rule_159a590c7b2c8d95),
rel_RuleMatch_8974a5cadf2d4779(&rel_RuleMatch_8974a5cadf2d4779){
}

void Stratum_RuleMatch_4394a245605301ce::run([[maybe_unused]] const std::vector<RamDomain>& args,[[maybe_unused]] std::vector<RamDomain>& ret){
signalHandler->setMsg(R"_(RuleMatch(name,location) :- 
   Rule(name,node),
   AST_NodeLocation(node,location).
in file rules.dl [8:1-8:81])_");
if(!(rel_Rule_159a590c7b2c8d95->empty()) && !(rel_AST_NodeLocation_5f3f38ee7a82c12a->empty())) {
[&](){
CREATE_OP_CONTEXT(rel_AST_NodeLocation_5f3f38ee7a82c12a_op_ctxt,rel_AST_NodeLocation_5f3f38ee7a82c12a->createContext());
CREATE_OP_CONTEXT(rel_Rule_159a590c7b2c8d95_op_ctxt,rel_Rule_159a590c7b2c8d95->createContext());
CREATE_OP_CONTEXT(rel_RuleMatch_8974a5cadf2d4779_op_ctxt,rel_RuleMatch_8974a5cadf2d4779->createContext());
for(const auto& env0 : *rel_Rule_159a590c7b2c8d95) {
auto range = rel_AST_NodeLocation_5f3f38ee7a82c12a->lowerUpperRange_10(Tuple<RamDomain,2>{{ramBitCast(env0[1]), ramBitCast<RamDomain>(MIN_RAM_SIGNED)}},Tuple<RamDomain,2>{{ramBitCast(env0[1]), ramBitCast<RamDomain>(MAX_RAM_SIGNED)}},READ_OP_CONTEXT(rel_AST_NodeLocation_5f3f38ee7a82c12a_op_ctxt));
for(const auto& env1 : range) {
Tuple<RamDomain,2> tuple{{ramBitCast(env0[0]),ramBitCast(env1[1])}};
rel_RuleMatch_8974a5cadf2d4779->insert(tuple,READ_OP_CONTEXT(rel_RuleMatch_8974a5cadf2d4779_op_ctxt));
}
}
}
();}
if (performIO) {
try {std::map<std::string, std::string> directiveMap({{"IO","file"},{"attributeNames","name\tlocation"},{"auxArity","0"},{"name","RuleMatch"},{"operation","output"},{"output-dir","."},{"params","{\"records\": {\"AST_Location\": {\"arity\": 5, \"params\": [\"startByte\", \"startLine\", \"startColumn\", \"endLine\", \"endColumn\"]}, \"AST_Node\": {\"arity\": 2, \"params\": [\"file\", \"nodeId\"]}}, \"relation\": {\"arity\": 2, \"params\": [\"name\", \"location\"]}}"},{"types","{\"ADTs\": {}, \"records\": {\"r:AST_Location\": {\"arity\": 5, \"types\": [\"i:AST_FilePosition\", \"i:AST_FilePosition\", \"i:AST_FilePosition\", \"i:AST_FilePosition\", \"i:AST_FilePosition\"]}, \"r:AST_Node\": {\"arity\": 2, \"types\": [\"i:AST_FileId\", \"i:AST_NodeId\"]}}, \"relation\": {\"arity\": 2, \"types\": [\"s:RuleName\", \"r:AST_Location\"]}}"}});
if (outputDirectory == "-"){directiveMap["IO"] = "stdout"; directiveMap["headers"] = "true";}
else if (!outputDirectory.empty()) {directiveMap["output-dir"] = outputDirectory;}
IOSystem::getInstance().getWriter(directiveMap, symTable, recordTable)->writeAll(*rel_RuleMatch_8974a5cadf2d4779);
} catch (std::exception& e) {std::cerr << e.what();exit(1);}
}
if (pruneImdtRels) rel_AST_NodeLocation_5f3f38ee7a82c12a->purge();
if (pruneImdtRels) rel_Rule_159a590c7b2c8d95->purge();
}

} // namespace  souffle

namespace  souffle {
using namespace souffle;
class Sf_generated: public SouffleProgram {
public:
 Sf_generated();
 ~Sf_generated();
void run();
void runAll(std::string inputDirectoryArg = "",std::string outputDirectoryArg = "",bool performIOArg = true,bool pruneImdtRelsArg = true);
void printAll([[maybe_unused]] std::string outputDirectoryArg = "");
void loadAll([[maybe_unused]] std::string inputDirectoryArg = "");
void dumpInputs();
void dumpOutputs();
SymbolTable& getSymbolTable();
RecordTable& getRecordTable();
void setNumThreads(std::size_t numThreadsValue);
void executeSubroutine(std::string name,const std::vector<RamDomain>& args,std::vector<RamDomain>& ret);
private:
void runFunction(std::string inputDirectoryArg,std::string outputDirectoryArg,bool performIOArg,bool pruneImdtRelsArg);
SymbolTableImpl symTable;
SpecializedRecordTable<0> recordTable;
ConcurrentCache<std::string,std::regex> regexCache;
Own<t_btree_ii__0_1__11::Type> rel_AST_NodeContent_b2f3666572e60754;
souffle::RelationWrapper<t_btree_ii__0_1__11::Type> wrapper_rel_AST_NodeContent_b2f3666572e60754;
Own<t_btree_iii__0_2_1__101__111::Type> rel_AST_NodeField_ca02670731ce3c99;
souffle::RelationWrapper<t_btree_iii__0_2_1__101__111::Type> wrapper_rel_AST_NodeField_ca02670731ce3c99;
Own<t_btree_ii__0_1__11__10::Type> rel_AST_NodeLocation_5f3f38ee7a82c12a;
souffle::RelationWrapper<t_btree_ii__0_1__11__10::Type> wrapper_rel_AST_NodeLocation_5f3f38ee7a82c12a;
Own<t_btree_ii__1_0__11__01::Type> rel_AST_NodeType_b38285ae9991409e;
souffle::RelationWrapper<t_btree_ii__1_0__11__01::Type> wrapper_rel_AST_NodeType_b38285ae9991409e;
Own<t_btree_iui__0_1_2__110__111::Type> rel_AST_ParentChild_be6259205eb66578;
souffle::RelationWrapper<t_btree_iui__0_1_2__110__111::Type> wrapper_rel_AST_ParentChild_be6259205eb66578;
Own<t_btree_ii__0_1__11::Type> rel_Rule_159a590c7b2c8d95;
souffle::RelationWrapper<t_btree_ii__0_1__11::Type> wrapper_rel_Rule_159a590c7b2c8d95;
Own<t_btree_ii__0_1__11::Type> rel_RuleMatch_8974a5cadf2d4779;
souffle::RelationWrapper<t_btree_ii__0_1__11::Type> wrapper_rel_RuleMatch_8974a5cadf2d4779;
Stratum_AST_NodeContent_fd51b4bf60caba3f stratum_AST_NodeContent_f29a0e907561c50c;
Stratum_AST_NodeField_cc21295739297165 stratum_AST_NodeField_b63714d335ba8b2e;
Stratum_AST_NodeLocation_89d765aa14237a09 stratum_AST_NodeLocation_b120359576603175;
Stratum_AST_NodeType_400775685fb3e630 stratum_AST_NodeType_f0647903c04c5ec5;
Stratum_AST_ParentChild_798cb83c96de8e4d stratum_AST_ParentChild_381ffe0668807a9c;
Stratum_Rule_0e6b7aa9ece342e5 stratum_Rule_b7f309df7d140ebc;
Stratum_RuleMatch_4394a245605301ce stratum_RuleMatch_a3becf205b4cd965;
std::string inputDirectory;
std::string outputDirectory;
SignalHandler* signalHandler{SignalHandler::instance()};
std::atomic<RamDomain> ctr{};
std::atomic<std::size_t> iter{};
};
} // namespace  souffle
namespace  souffle {
using namespace souffle;
 Sf_generated::Sf_generated():
symTable({
	R"_(call)_",
	R"_(receiver)_",
	R"_(block)_",
	R"_(method)_",
	R"_(CREATE)_",
	R"_(constant)_",
	R"_(arguments)_",
	R"_(argument_list)_",
	R"_(TABLE)_",
	R"_(public)_",
	R"_(identifier)_",
	R"_(.)_",
	R"_()_",
	R"_(global_variable)_",
	R"_(parenthesized_statements)_",
	R"_(()_",
	R"_(ERROR)_",
	R"_(<)_",
	R"_(>)_",
	R"_())_",
	R"_(sql_lang_create_table)_",
}),
recordTable(),
regexCache(),
rel_AST_NodeContent_b2f3666572e60754(mk<t_btree_ii__0_1__11::Type>()),
wrapper_rel_AST_NodeContent_b2f3666572e60754(0, *rel_AST_NodeContent_b2f3666572e60754, *this, "AST_NodeContent", std::array<const char *,2>{{"r:AST_Node","s:AST_Content"}}, std::array<const char *,2>{{"node","content"}}, 0),
rel_AST_NodeField_ca02670731ce3c99(mk<t_btree_iii__0_2_1__101__111::Type>()),
wrapper_rel_AST_NodeField_ca02670731ce3c99(1, *rel_AST_NodeField_ca02670731ce3c99, *this, "AST_NodeField", std::array<const char *,3>{{"r:AST_Node","r:AST_Node","s:AST_Field"}}, std::array<const char *,3>{{"parent","child","field"}}, 0),
rel_AST_NodeLocation_5f3f38ee7a82c12a(mk<t_btree_ii__0_1__11__10::Type>()),
wrapper_rel_AST_NodeLocation_5f3f38ee7a82c12a(2, *rel_AST_NodeLocation_5f3f38ee7a82c12a, *this, "AST_NodeLocation", std::array<const char *,2>{{"r:AST_Node","r:AST_Location"}}, std::array<const char *,2>{{"node","location"}}, 0),
rel_AST_NodeType_b38285ae9991409e(mk<t_btree_ii__1_0__11__01::Type>()),
wrapper_rel_AST_NodeType_b38285ae9991409e(3, *rel_AST_NodeType_b38285ae9991409e, *this, "AST_NodeType", std::array<const char *,2>{{"r:AST_Node","s:AST_Type"}}, std::array<const char *,2>{{"node","type"}}, 0),
rel_AST_ParentChild_be6259205eb66578(mk<t_btree_iui__0_1_2__110__111::Type>()),
wrapper_rel_AST_ParentChild_be6259205eb66578(4, *rel_AST_ParentChild_be6259205eb66578, *this, "AST_ParentChild", std::array<const char *,3>{{"r:AST_Node","u:unsigned","r:AST_Node"}}, std::array<const char *,3>{{"parent","index","child"}}, 0),
rel_Rule_159a590c7b2c8d95(mk<t_btree_ii__0_1__11::Type>()),
wrapper_rel_Rule_159a590c7b2c8d95(5, *rel_Rule_159a590c7b2c8d95, *this, "Rule", std::array<const char *,2>{{"s:RuleName","r:AST_Node"}}, std::array<const char *,2>{{"rule","node"}}, 0),
rel_RuleMatch_8974a5cadf2d4779(mk<t_btree_ii__0_1__11::Type>()),
wrapper_rel_RuleMatch_8974a5cadf2d4779(6, *rel_RuleMatch_8974a5cadf2d4779, *this, "RuleMatch", std::array<const char *,2>{{"s:RuleName","r:AST_Location"}}, std::array<const char *,2>{{"name","location"}}, 0),
stratum_AST_NodeContent_f29a0e907561c50c(symTable,recordTable,regexCache,pruneImdtRels,performIO,signalHandler,iter,ctr,inputDirectory,outputDirectory,*rel_AST_NodeContent_b2f3666572e60754),
stratum_AST_NodeField_b63714d335ba8b2e(symTable,recordTable,regexCache,pruneImdtRels,performIO,signalHandler,iter,ctr,inputDirectory,outputDirectory,*rel_AST_NodeField_ca02670731ce3c99),
stratum_AST_NodeLocation_b120359576603175(symTable,recordTable,regexCache,pruneImdtRels,performIO,signalHandler,iter,ctr,inputDirectory,outputDirectory,*rel_AST_NodeLocation_5f3f38ee7a82c12a),
stratum_AST_NodeType_f0647903c04c5ec5(symTable,recordTable,regexCache,pruneImdtRels,performIO,signalHandler,iter,ctr,inputDirectory,outputDirectory,*rel_AST_NodeType_b38285ae9991409e),
stratum_AST_ParentChild_381ffe0668807a9c(symTable,recordTable,regexCache,pruneImdtRels,performIO,signalHandler,iter,ctr,inputDirectory,outputDirectory,*rel_AST_ParentChild_be6259205eb66578),
stratum_Rule_b7f309df7d140ebc(symTable,recordTable,regexCache,pruneImdtRels,performIO,signalHandler,iter,ctr,inputDirectory,outputDirectory,*rel_AST_NodeContent_b2f3666572e60754,*rel_AST_NodeField_ca02670731ce3c99,*rel_AST_NodeType_b38285ae9991409e,*rel_AST_ParentChild_be6259205eb66578,*rel_Rule_159a590c7b2c8d95),
stratum_RuleMatch_a3becf205b4cd965(symTable,recordTable,regexCache,pruneImdtRels,performIO,signalHandler,iter,ctr,inputDirectory,outputDirectory,*rel_AST_NodeLocation_5f3f38ee7a82c12a,*rel_Rule_159a590c7b2c8d95,*rel_RuleMatch_8974a5cadf2d4779){
addRelation("AST_NodeContent", wrapper_rel_AST_NodeContent_b2f3666572e60754, true, false);
addRelation("AST_NodeField", wrapper_rel_AST_NodeField_ca02670731ce3c99, true, false);
addRelation("AST_NodeLocation", wrapper_rel_AST_NodeLocation_5f3f38ee7a82c12a, true, false);
addRelation("AST_NodeType", wrapper_rel_AST_NodeType_b38285ae9991409e, true, false);
addRelation("AST_ParentChild", wrapper_rel_AST_ParentChild_be6259205eb66578, true, false);
addRelation("Rule", wrapper_rel_Rule_159a590c7b2c8d95, false, false);
addRelation("RuleMatch", wrapper_rel_RuleMatch_8974a5cadf2d4779, false, true);
}

 Sf_generated::~Sf_generated(){
}

void Sf_generated::runFunction(std::string inputDirectoryArg,std::string outputDirectoryArg,bool performIOArg,bool pruneImdtRelsArg){

    this->inputDirectory  = std::move(inputDirectoryArg);
    this->outputDirectory = std::move(outputDirectoryArg);
    this->performIO       = performIOArg;
    this->pruneImdtRels   = pruneImdtRelsArg;

    // set default threads (in embedded mode)
    // if this is not set, and omp is used, the default omp setting of number of cores is used.
#if defined(_OPENMP)
    if (0 < getNumThreads()) { omp_set_num_threads(static_cast<int>(getNumThreads())); }
#endif

    signalHandler->set();
// -- query evaluation --
{
 std::vector<RamDomain> args, ret;
stratum_AST_NodeContent_f29a0e907561c50c.run(args, ret);
}
{
 std::vector<RamDomain> args, ret;
stratum_AST_NodeField_b63714d335ba8b2e.run(args, ret);
}
{
 std::vector<RamDomain> args, ret;
stratum_AST_NodeLocation_b120359576603175.run(args, ret);
}
{
 std::vector<RamDomain> args, ret;
stratum_AST_NodeType_f0647903c04c5ec5.run(args, ret);
}
{
 std::vector<RamDomain> args, ret;
stratum_AST_ParentChild_381ffe0668807a9c.run(args, ret);
}
{
 std::vector<RamDomain> args, ret;
stratum_Rule_b7f309df7d140ebc.run(args, ret);
}
{
 std::vector<RamDomain> args, ret;
stratum_RuleMatch_a3becf205b4cd965.run(args, ret);
}

// -- relation hint statistics --
signalHandler->reset();
}

void Sf_generated::run(){
runFunction("", "", false, false);
}

void Sf_generated::runAll(std::string inputDirectoryArg,std::string outputDirectoryArg,bool performIOArg,bool pruneImdtRelsArg){
runFunction(inputDirectoryArg, outputDirectoryArg, performIOArg, pruneImdtRelsArg);
}

void Sf_generated::printAll([[maybe_unused]] std::string outputDirectoryArg){
try {std::map<std::string, std::string> directiveMap({{"IO","file"},{"attributeNames","name\tlocation"},{"auxArity","0"},{"name","RuleMatch"},{"operation","output"},{"output-dir","."},{"params","{\"records\": {\"AST_Location\": {\"arity\": 5, \"params\": [\"startByte\", \"startLine\", \"startColumn\", \"endLine\", \"endColumn\"]}, \"AST_Node\": {\"arity\": 2, \"params\": [\"file\", \"nodeId\"]}}, \"relation\": {\"arity\": 2, \"params\": [\"name\", \"location\"]}}"},{"types","{\"ADTs\": {}, \"records\": {\"r:AST_Location\": {\"arity\": 5, \"types\": [\"i:AST_FilePosition\", \"i:AST_FilePosition\", \"i:AST_FilePosition\", \"i:AST_FilePosition\", \"i:AST_FilePosition\"]}, \"r:AST_Node\": {\"arity\": 2, \"types\": [\"i:AST_FileId\", \"i:AST_NodeId\"]}}, \"relation\": {\"arity\": 2, \"types\": [\"s:RuleName\", \"r:AST_Location\"]}}"}});
if (!outputDirectoryArg.empty()) {directiveMap["output-dir"] = outputDirectoryArg;}
IOSystem::getInstance().getWriter(directiveMap, symTable, recordTable)->writeAll(*rel_RuleMatch_8974a5cadf2d4779);
} catch (std::exception& e) {std::cerr << e.what();exit(1);}
}

void Sf_generated::loadAll([[maybe_unused]] std::string inputDirectoryArg){
try {std::map<std::string, std::string> directiveMap({{"IO","file"},{"attributeNames","node\tcontent"},{"auxArity","0"},{"fact-dir","."},{"name","AST_NodeContent"},{"operation","input"},{"params","{\"records\": {\"AST_Location\": {\"arity\": 5, \"params\": [\"startByte\", \"startLine\", \"startColumn\", \"endLine\", \"endColumn\"]}, \"AST_Node\": {\"arity\": 2, \"params\": [\"file\", \"nodeId\"]}}, \"relation\": {\"arity\": 2, \"params\": [\"node\", \"content\"]}}"},{"types","{\"ADTs\": {}, \"records\": {\"r:AST_Location\": {\"arity\": 5, \"types\": [\"i:AST_FilePosition\", \"i:AST_FilePosition\", \"i:AST_FilePosition\", \"i:AST_FilePosition\", \"i:AST_FilePosition\"]}, \"r:AST_Node\": {\"arity\": 2, \"types\": [\"i:AST_FileId\", \"i:AST_NodeId\"]}}, \"relation\": {\"arity\": 2, \"types\": [\"r:AST_Node\", \"s:AST_Content\"]}}"}});
if (!inputDirectoryArg.empty()) {directiveMap["fact-dir"] = inputDirectoryArg;}
IOSystem::getInstance().getReader(directiveMap, symTable, recordTable)->readAll(*rel_AST_NodeContent_b2f3666572e60754);
} catch (std::exception& e) {std::cerr << "Error loading AST_NodeContent data: " << e.what() << '\n';
exit(1);
}
try {std::map<std::string, std::string> directiveMap({{"IO","file"},{"attributeNames","parent\tchild\tfield"},{"auxArity","0"},{"fact-dir","."},{"name","AST_NodeField"},{"operation","input"},{"params","{\"records\": {\"AST_Location\": {\"arity\": 5, \"params\": [\"startByte\", \"startLine\", \"startColumn\", \"endLine\", \"endColumn\"]}, \"AST_Node\": {\"arity\": 2, \"params\": [\"file\", \"nodeId\"]}}, \"relation\": {\"arity\": 3, \"params\": [\"parent\", \"child\", \"field\"]}}"},{"types","{\"ADTs\": {}, \"records\": {\"r:AST_Location\": {\"arity\": 5, \"types\": [\"i:AST_FilePosition\", \"i:AST_FilePosition\", \"i:AST_FilePosition\", \"i:AST_FilePosition\", \"i:AST_FilePosition\"]}, \"r:AST_Node\": {\"arity\": 2, \"types\": [\"i:AST_FileId\", \"i:AST_NodeId\"]}}, \"relation\": {\"arity\": 3, \"types\": [\"r:AST_Node\", \"r:AST_Node\", \"s:AST_Field\"]}}"}});
if (!inputDirectoryArg.empty()) {directiveMap["fact-dir"] = inputDirectoryArg;}
IOSystem::getInstance().getReader(directiveMap, symTable, recordTable)->readAll(*rel_AST_NodeField_ca02670731ce3c99);
} catch (std::exception& e) {std::cerr << "Error loading AST_NodeField data: " << e.what() << '\n';
exit(1);
}
try {std::map<std::string, std::string> directiveMap({{"IO","file"},{"attributeNames","node\tlocation"},{"auxArity","0"},{"fact-dir","."},{"name","AST_NodeLocation"},{"operation","input"},{"params","{\"records\": {\"AST_Location\": {\"arity\": 5, \"params\": [\"startByte\", \"startLine\", \"startColumn\", \"endLine\", \"endColumn\"]}, \"AST_Node\": {\"arity\": 2, \"params\": [\"file\", \"nodeId\"]}}, \"relation\": {\"arity\": 2, \"params\": [\"node\", \"location\"]}}"},{"types","{\"ADTs\": {}, \"records\": {\"r:AST_Location\": {\"arity\": 5, \"types\": [\"i:AST_FilePosition\", \"i:AST_FilePosition\", \"i:AST_FilePosition\", \"i:AST_FilePosition\", \"i:AST_FilePosition\"]}, \"r:AST_Node\": {\"arity\": 2, \"types\": [\"i:AST_FileId\", \"i:AST_NodeId\"]}}, \"relation\": {\"arity\": 2, \"types\": [\"r:AST_Node\", \"r:AST_Location\"]}}"}});
if (!inputDirectoryArg.empty()) {directiveMap["fact-dir"] = inputDirectoryArg;}
IOSystem::getInstance().getReader(directiveMap, symTable, recordTable)->readAll(*rel_AST_NodeLocation_5f3f38ee7a82c12a);
} catch (std::exception& e) {std::cerr << "Error loading AST_NodeLocation data: " << e.what() << '\n';
exit(1);
}
try {std::map<std::string, std::string> directiveMap({{"IO","file"},{"attributeNames","node\ttype"},{"auxArity","0"},{"fact-dir","."},{"name","AST_NodeType"},{"operation","input"},{"params","{\"records\": {\"AST_Location\": {\"arity\": 5, \"params\": [\"startByte\", \"startLine\", \"startColumn\", \"endLine\", \"endColumn\"]}, \"AST_Node\": {\"arity\": 2, \"params\": [\"file\", \"nodeId\"]}}, \"relation\": {\"arity\": 2, \"params\": [\"node\", \"type\"]}}"},{"types","{\"ADTs\": {}, \"records\": {\"r:AST_Location\": {\"arity\": 5, \"types\": [\"i:AST_FilePosition\", \"i:AST_FilePosition\", \"i:AST_FilePosition\", \"i:AST_FilePosition\", \"i:AST_FilePosition\"]}, \"r:AST_Node\": {\"arity\": 2, \"types\": [\"i:AST_FileId\", \"i:AST_NodeId\"]}}, \"relation\": {\"arity\": 2, \"types\": [\"r:AST_Node\", \"s:AST_Type\"]}}"}});
if (!inputDirectoryArg.empty()) {directiveMap["fact-dir"] = inputDirectoryArg;}
IOSystem::getInstance().getReader(directiveMap, symTable, recordTable)->readAll(*rel_AST_NodeType_b38285ae9991409e);
} catch (std::exception& e) {std::cerr << "Error loading AST_NodeType data: " << e.what() << '\n';
exit(1);
}
try {std::map<std::string, std::string> directiveMap({{"IO","file"},{"attributeNames","parent\tindex\tchild"},{"auxArity","0"},{"fact-dir","."},{"name","AST_ParentChild"},{"operation","input"},{"params","{\"records\": {\"AST_Location\": {\"arity\": 5, \"params\": [\"startByte\", \"startLine\", \"startColumn\", \"endLine\", \"endColumn\"]}, \"AST_Node\": {\"arity\": 2, \"params\": [\"file\", \"nodeId\"]}}, \"relation\": {\"arity\": 3, \"params\": [\"parent\", \"index\", \"child\"]}}"},{"types","{\"ADTs\": {}, \"records\": {\"r:AST_Location\": {\"arity\": 5, \"types\": [\"i:AST_FilePosition\", \"i:AST_FilePosition\", \"i:AST_FilePosition\", \"i:AST_FilePosition\", \"i:AST_FilePosition\"]}, \"r:AST_Node\": {\"arity\": 2, \"types\": [\"i:AST_FileId\", \"i:AST_NodeId\"]}}, \"relation\": {\"arity\": 3, \"types\": [\"r:AST_Node\", \"u:unsigned\", \"r:AST_Node\"]}}"}});
if (!inputDirectoryArg.empty()) {directiveMap["fact-dir"] = inputDirectoryArg;}
IOSystem::getInstance().getReader(directiveMap, symTable, recordTable)->readAll(*rel_AST_ParentChild_be6259205eb66578);
} catch (std::exception& e) {std::cerr << "Error loading AST_ParentChild data: " << e.what() << '\n';
exit(1);
}
}

void Sf_generated::dumpInputs(){
try {std::map<std::string, std::string> rwOperation;
rwOperation["IO"] = "stdout";
rwOperation["name"] = "AST_NodeContent";
rwOperation["types"] = "{\"relation\": {\"arity\": 2, \"auxArity\": 0, \"types\": [\"r:AST_Node\", \"s:AST_Content\"]}}";
IOSystem::getInstance().getWriter(rwOperation, symTable, recordTable)->writeAll(*rel_AST_NodeContent_b2f3666572e60754);
} catch (std::exception& e) {std::cerr << e.what();exit(1);}
try {std::map<std::string, std::string> rwOperation;
rwOperation["IO"] = "stdout";
rwOperation["name"] = "AST_NodeField";
rwOperation["types"] = "{\"relation\": {\"arity\": 3, \"auxArity\": 0, \"types\": [\"r:AST_Node\", \"r:AST_Node\", \"s:AST_Field\"]}}";
IOSystem::getInstance().getWriter(rwOperation, symTable, recordTable)->writeAll(*rel_AST_NodeField_ca02670731ce3c99);
} catch (std::exception& e) {std::cerr << e.what();exit(1);}
try {std::map<std::string, std::string> rwOperation;
rwOperation["IO"] = "stdout";
rwOperation["name"] = "AST_NodeLocation";
rwOperation["types"] = "{\"relation\": {\"arity\": 2, \"auxArity\": 0, \"types\": [\"r:AST_Node\", \"r:AST_Location\"]}}";
IOSystem::getInstance().getWriter(rwOperation, symTable, recordTable)->writeAll(*rel_AST_NodeLocation_5f3f38ee7a82c12a);
} catch (std::exception& e) {std::cerr << e.what();exit(1);}
try {std::map<std::string, std::string> rwOperation;
rwOperation["IO"] = "stdout";
rwOperation["name"] = "AST_NodeType";
rwOperation["types"] = "{\"relation\": {\"arity\": 2, \"auxArity\": 0, \"types\": [\"r:AST_Node\", \"s:AST_Type\"]}}";
IOSystem::getInstance().getWriter(rwOperation, symTable, recordTable)->writeAll(*rel_AST_NodeType_b38285ae9991409e);
} catch (std::exception& e) {std::cerr << e.what();exit(1);}
try {std::map<std::string, std::string> rwOperation;
rwOperation["IO"] = "stdout";
rwOperation["name"] = "AST_ParentChild";
rwOperation["types"] = "{\"relation\": {\"arity\": 3, \"auxArity\": 0, \"types\": [\"r:AST_Node\", \"u:unsigned\", \"r:AST_Node\"]}}";
IOSystem::getInstance().getWriter(rwOperation, symTable, recordTable)->writeAll(*rel_AST_ParentChild_be6259205eb66578);
} catch (std::exception& e) {std::cerr << e.what();exit(1);}
}

void Sf_generated::dumpOutputs(){
try {std::map<std::string, std::string> rwOperation;
rwOperation["IO"] = "stdout";
rwOperation["name"] = "RuleMatch";
rwOperation["types"] = "{\"relation\": {\"arity\": 2, \"auxArity\": 0, \"types\": [\"s:RuleName\", \"r:AST_Location\"]}}";
IOSystem::getInstance().getWriter(rwOperation, symTable, recordTable)->writeAll(*rel_RuleMatch_8974a5cadf2d4779);
} catch (std::exception& e) {std::cerr << e.what();exit(1);}
}

SymbolTable& Sf_generated::getSymbolTable(){
return symTable;
}

RecordTable& Sf_generated::getRecordTable(){
return recordTable;
}

void Sf_generated::setNumThreads(std::size_t numThreadsValue){
SouffleProgram::setNumThreads(numThreadsValue);
symTable.setNumLanes(getNumThreads());
recordTable.setNumLanes(getNumThreads());
regexCache.setNumLanes(getNumThreads());
}

void Sf_generated::executeSubroutine(std::string name,const std::vector<RamDomain>& args,std::vector<RamDomain>& ret){
if (name == "AST_NodeContent") {
stratum_AST_NodeContent_f29a0e907561c50c.run(args, ret);
return;}
if (name == "AST_NodeField") {
stratum_AST_NodeField_b63714d335ba8b2e.run(args, ret);
return;}
if (name == "AST_NodeLocation") {
stratum_AST_NodeLocation_b120359576603175.run(args, ret);
return;}
if (name == "AST_NodeType") {
stratum_AST_NodeType_f0647903c04c5ec5.run(args, ret);
return;}
if (name == "AST_ParentChild") {
stratum_AST_ParentChild_381ffe0668807a9c.run(args, ret);
return;}
if (name == "Rule") {
stratum_Rule_b7f309df7d140ebc.run(args, ret);
return;}
if (name == "RuleMatch") {
stratum_RuleMatch_a3becf205b4cd965.run(args, ret);
return;}
fatal(("unknown subroutine " + name).c_str());
}

} // namespace  souffle
namespace souffle {
SouffleProgram *newInstance_generated(){return new  souffle::Sf_generated;}
SymbolTable *getST_generated(SouffleProgram *p){return &reinterpret_cast<souffle::Sf_generated*>(p)->getSymbolTable();}
} // namespace souffle

#ifndef __EMBEDDED_SOUFFLE__
#include "souffle/CompiledOptions.h"
int main(int argc, char** argv)
{
try{
souffle::CmdOptions opt(R"(souffle/rules.dl)",
R"()",
R"()",
false,
R"()",
1);
if (!opt.parse(argc,argv)) return 1;
souffle::Sf_generated obj;
#if defined(_OPENMP) 
obj.setNumThreads(opt.getNumJobs());

#endif
obj.runAll(opt.getInputFileDir(), opt.getOutputFileDir());
return 0;
} catch(std::exception &e) { souffle::SignalHandler::instance()->error(e.what());}
}
#endif

namespace  souffle {
using namespace souffle;
class factory_Sf_generated: souffle::ProgramFactory {
public:
souffle::SouffleProgram* newInstance();
 factory_Sf_generated();
private:
};
} // namespace  souffle
namespace  souffle {
using namespace souffle;
souffle::SouffleProgram* factory_Sf_generated::newInstance(){
return new  souffle::Sf_generated();
}

 factory_Sf_generated::factory_Sf_generated():
souffle::ProgramFactory("generated"){
}

} // namespace  souffle
namespace souffle {

#ifdef __EMBEDDED_SOUFFLE__
extern "C" {
souffle::factory_Sf_generated __factory_Sf_generated_instance;
}
#endif
} // namespace souffle


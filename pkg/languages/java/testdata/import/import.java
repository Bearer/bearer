import foo.Import;
import foo.Import2.*;
import static foo.Import3;

class A {
    public void exec() {
        sink(Import);
        sink(Import2); // no match
        sink(Import3); // no match
    }
}
<?php
// class test
class Foo {
    public $aMemberVar = 'aMemberVar Member Variable';
   	static $test;
   
    function aMemberFunc() {
        print 'Inside `aMemberFunc()`';

        $this->test->subtest;
        $this->test2->subtest2;
    }
}


// property -> acess test
$parent->child;
$parent->childMethod();
$parent->$child->grandchild;
$parent->child->grandchildMethod()

$parent->childMethod()->subParent->subChildMethod();
$parent->child->grandChildMethod()->subParent->subChild->subGrandChildMethod();


// scoping tests
$test1->test();
$test1->test();

function () {
    $test1->test();
    $test1->test();
}

function Test(){
    $test1->test();
    $test1->test();
}

class Test {
    static function method() {
        $test1->test();
        $test1->test();
    }
}

?>
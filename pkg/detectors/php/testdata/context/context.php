<?php

$testglobal = "global"; // supported
$test = "test global"; // supported

$anony = function ($a , $b){
  $test = 'test anon function'; // supported
  "https://test.api.com/".$test."/".$testglobal;
  return $a < $b;
 };
 
$anony = new class {
    public $test = 'test anon class property'; // supported
    public $testmember = "test anon class property"; // supported
    public function hello($name) {
        $test = 'test anon class member function'; // supported
        "https://test.api.com/".$test."/".$testmember.$testglobal; // supported
        "https://test.api.com/".$this->$test."/".$this->$testmember; // not supported
        return "Hello $name";
    }
};

class Foo {
    public $test = 'test class member'; // supported
    public $testmember = "test class property"; // supported

    function aMemberFunc() {
        $test = 'test class member function'; // supported
        "https://test.api.com/".$test."/".$testmember.$testglobal; // supported
        "https://test.api.com/".$this->$test."/".$this->$testmember; // not supported

        print 'Inside `aMemberFunc()`';
    }

    public function __construct(int $id, ?string $name)
    {
        $test = 'test class member constructor'; // supported
        "https://test.api.com/".$test."/".$testmember."/".$testglobal; // supported
        "https://test.api.com/".$this->$test."/".$this->$testmember; // not supported
        $this->id = $id;
        $this->name = $name;
    }
}
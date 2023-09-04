<?php
class Greet {
    const Greeting = "Hello World";

    public static function main($args)
    {
        $s = self::Greeting . "!";
        $s .= "!!";

        $s2 = "hey ";
        $s2 .= $args[0];
        $s2 .= " there";
    }
}
?>
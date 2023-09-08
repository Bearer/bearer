<?php
scopeCursor($_GET["oops"]);
scopeCursor(x . $_GET["oops"]);
scopeCursor(x ? $_GET["oops"] : y);
scopeCursor($_GET["ok"] ? x : y);
scopeCursor($_GET["oops"] ?: y);

scopeNested($_GET["oops"]);
scopeNested(x . $_GET["oops"]);
scopeNested(x ? $_GET["oops"] : y);
scopeNested($_GET["oops"] ? x : y);
scopeNested($_GET["oops"] ?: y);

scopeResult($_GET["oops"]);
scopeResult(x . $_GET["oops"]);
scopeResult(x ? $_GET["oops"] : y);
scopeResult($_GET["ok"] ? x : y);
scopeResult($_GET["oops"] ?: y);
?>
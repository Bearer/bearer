<?php

$orderServiceUrl = $_ENV["ORDER_SERVICE_URL"];
$userServiceHost = $_ENV['USER_SERVICE_HOST'];
$accountId = $_ENV['ACCOUNT_ID'];
$other = $_POST["IGNORE_ME_HOST"];

$concat = $_ENV["CUSTOMERS_URL"] . "/path";
$interpolation = "{$_ENV["CUSTOMERS_HOST"]}:{$port}/path";

$x = array("ignored.domain.com" => "abc");
$y = $x["ignored.domain.com"];
?>

<html> 

<head>
</head>
<body>
    <script>var url = "https://inline.domain.com"</script>
	<script>
		var url = "https://api1.domain.com"
	</script>
	<script type="text/javascript">
		var url = "https://api2.domain.com"
	</script>
	<script type="text/json">
		{"url":"https://ignored-api.domain.com"}
	</script>
	<div class="links">
		<a href="https://ignored-domain.com">Foo</a>
		<a href="https://ignored-domain2.com">Bar</a>
	</div>


<?php
    $test= "test";
	$personalDataUrl = $_ENV["PERSONAL_DATA_URL"];
?>

    <script>
		var url = "https://api3.domain.com"
    </script>
</body>
</html>

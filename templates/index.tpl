<html>
	<head>
		<title>GChat Example</title>
		<script src="http://code.jquery.com/jquery-latest.min.js "></script>
		<script src="/static/js/golem.js"></script>
		<script src="/static/js/scripts.js"></script>
	</head>
	<body>
		<div id="uname">{{.}}</div>
		<a href="/logout">Logout</a>
		<div>
			<input id="msg" type="text" ><button id="send">Send</button>
			<ul class="users"></ul>
			<ul class="messages"></ul>
		</div>
	</body>
</html>
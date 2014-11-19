$(document).ready(function() {

	var uname = $('#uname').text();
	var conn = new golem.Connection("ws://127.0.0.1:3333/ws?uname=" + uname, true);

	conn.on("open", function() {
		console.log("GChat: Connected.");
	});

	conn.on("close", function() {
		console.log("GChat: Disconnected.");
	});

	conn.on("join", function(data) {
		console.log("GChat: Join " + data.name);
		if (data.name !== uname) {
			$(".messages").append("<li class='item joined'>"+data.name+" joins</li>");
			$(".users").append("<li class='item user "+data.name+"'>"+data.name+"</li>");
		} else {
			$(data.messages).each(function(indx, element){
				$(".messages").append("<li class='item msg'>"+element.uname+": "+element.msg+"</li>");
			});
			console.log(data.users);
			$(data.users).each(function(indx, element){
				$(".users").append("<li class='item user "+element+"'>"+element+"</li>");
			});
		}

	});

	conn.on("leave", function(name) {
		console.log("GChat: Leave " + name);
		if (name !== uname) {
			$(".messages").append("<li class='item leave'>"+name+" leave chat</li>");
		};
		$(".users li.item.user."+name).remove();
	});

	conn.on("message", function(data) {
		$(".messages").append("<li class='item msg'>"+data.uname+": "+data.msg+"</li>");
		console.log(data);
	});

	var msg = $('#msg');
	$('button#send').click(function(){
		if (msg.val()) {
			conn.emit("msg", {uname: uname, msg: msg.val()});
			msg.val('');
		};
	});

});
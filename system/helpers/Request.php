<?php
class Request
{
	public static function method()
	{
		if (isset($_SERVER['HTTP_X_REQUESTED_WITH']) and $_SERVER['HTTP_X_REQUESTED_WITH'] == 'XMLHttpRequest') return 'AJAX';
		else if (!empty($_POST)) return 'POST';
		else return 'GET';
	}
}
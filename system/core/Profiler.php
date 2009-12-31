<?php

final class Profiler
{
	static function displayTrace($e)
	{
		echo '<style>h1,h2,h3,p,td{font-family:Verdana;font-weight:lighter;}</style>';
		echo '<p>Uncaught '.get_class($e).'</p><h1>'.$e->getMessage().'</h1>';

		$traces = $e->getTrace();
		if (count($traces) > 1) {
			echo '<p><b>Trace in execution order:</b></p>'.
				 '<pre style="font-family:Verdana; line-height: 20px">';

			$level = 0;
			foreach (array_reverse($traces) as $trace)
			{
				$level++;
				if (isset($trace['class'])) echo $trace['class'].'&rarr;';
				$args = array();
				if ( ! empty($trace['args']))
				{
					foreach ($trace['args'] as $arg)
					{
						if (is_null($arg)) $args[] = 'null';
						else if (is_array($arg)) $args[] = 'array['.sizeof($arg).']';
						else if (is_object($arg)) $args[] = get_class($arg).' Object';
						else if (is_bool($arg)) $args[] = $arg ? 'true' : 'false';
						else if (is_int($arg)) $args[] = $arg;
						else $args[] = "'".htmlspecialchars(substr($arg, 0, 64)).(strlen($arg) > 63 ? '&hellip;': '')."'";
					}
				}
				echo '<b>'.$trace['function'].'</b>('.implode(', ',$args).')  '
					.'on line <code>'.(isset($trace['line']) ? $trace['line'] : 'unknown').'</code> '
					.'in <code>'.(isset($trace['file']) ? $trace['file'] : 'unknown')."</code>\n"
					.str_repeat("\t", $level);
			}
			echo '</pre>';
		}
		echo "<p>Exception was thrown on line <code>"
			. $e->getLine() . "</code> in <code>"
			. $e->getFile() . "</code></p>";
	}

	static function display()
	{
		$status = array(
			'requested url'  => SE::getCurrentUrl(),
			'controller'     => SE::getController(),
			'action'         => SE::getAction(),
			'params'         => 'array('.join(', ',SE::getParams()).')',
			'request method' => Request::method()
		);
		self::displayTable($status, 'Dispatcher status');

		$markers = array();
		$old_mark = '';
		foreach(Benchmark::$mark as $mark => $time)
		{
			$markers[$mark] = Benchmark::time($old_mark, $mark);
			$old_mark = $mark;
		}
		self::displayTable($markers, 'Benchmark');

		if (!empty($_GET)) self::displayTable($_GET, 'GET');
		if (!empty($_POST)) self::displayTable($_POST, 'POST');
		if (!empty($_COOKIE)) self::displayTable($_COOKIE, 'COOKIE');
		self::displayTable($_SERVER, 'SERVER');
	}

	static function displayTable($array, $label, $key_label='Variable', $value_label='Value')
	{
		echo '<h2>'.$label.'</h2>'
			.'<table cellpadding="3" cellspacing="0" style="width: 800px; border: 1px solid #ccc">'
			.'<tr><td style="border-right: 1px solid #ccc; border-bottom: 1px solid #ccc;">'.$key_label.'</td>'
			.'<td style="border-bottom: 1px solid #ccc;">'.$value_label.'</td></tr>';

		foreach ($array as $key => $value)
		{
			if (is_null($value)) $value = 'null';
			else if (is_array($value)) $value = 'array['.sizeof($value).']';
			else if (is_object($value)) $value = get_class($value).' Object';
			else if (is_bool($value)) $value = $value ? 'true' : 'false';
			else if (is_int($value)) $value = $value;
			else $value = htmlspecialchars(substr($value, 0, 64)).(strlen($value) > 63 ? '&hellip;': '');
			echo '<tr><td><code>'.$key.'</code></td><td><code>'.$value.'</code></td></tr>';
		}
		echo '</table>';
	}
}
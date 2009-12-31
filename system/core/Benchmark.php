<?php
/**
 *  Class: Benchmark
 *
 *  This class enables you to mark points and calculate the time difference
 *  between them.  Memory consumption can also be displayed.
 */
final class Benchmark
{
	public static $marks = array();

	/*
	   Method: Start
	   Set a benchmark start point.
	*/
	public static function start($name)
	{
		if (!isset(self::$marks[$name]))
		{
			self::$marks[$name] = array
			(
				'start'        => microtime(true),
				'stop'         => false,
				'memory_start' => function_exists('memory_get_usage') ? memory_get_usage() : 0,
				'memory_stop'  => false
			);
		}
	}

	/*
	   Method: stop
	   Set a benchmark stop point.
	*/
	public static function stop($name)
	{
		if (isset(self::$marks[$name]))
		{
			self::$marks[$name]['stop'] = microtime(true);
			self::$marks[$name]['memory_stop'] = function_exists('memory_get_usage') ? memory_get_usage() : 0;
		}
	}

	/*
	   Method: get
	   Get the elapsed time between a start and stop of a mark name, TRUE for all.
	*/
	public static function get($name, $decimals = 4)
	{
		if ($name === true)
		{
			$times = array();

			foreach(array_keys(self::$marks) as $name)
				$times[$name] = self::get($name, $decimals);

			return $times;
		}

		if (!isset(self::$marks[$name]))
			return false;

		if (self::$marks[$name]['stop'] === false)
			self::stop($name);

		return array
		(
			'time'   => number_format(self::$marks[$name]['stop'] - self::$marks[$name]['start'], $decimals),
			'memory' => convert_size(self::$marks[$name]['memory_stop'] - self::$marks[$name]['memory_start'])
		);
	}

} // end Benchmark class

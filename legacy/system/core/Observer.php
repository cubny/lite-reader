<?php
/*
   Class: Observer

   This is the event observer class

   About: system.init
    Called just after including all files to the core framework

   About: system.page_not_found
    Called just before processing the 404.php default view.

   About: system.execute
    Controller locating and initialization. A controller object will be
    created and executed just after

   About: system.display
    Displays the output that the framework has generated.

   About: system.shutdown
    Last event to run, just before PHP starts to shut down.
*/

final class Observer
{
	private static $events = array(); // events callback

	/*
	   Method: observe
	   Attach a callback to an event queue.
	*/
	public static function observe($name, $callback)
	{
		if ( ! isset(self::$events[$name]))
			self::$events[$name] = array();

		self::$events[$name][] = $callback;
	}

	/*
	   Method: clear
	   Detach a callback to an event queue.
	*/
	public static function clear($name, $callback=false)
	{
		if ( ! $callback)
		{
			self::$events[$name] = array();
		}
		else if (isset(self::$events[$name]))
		{
			foreach (self::$events[$name] as $i => $event_callback)
			{
				if ($callback === $event_callback)
					unset(self::$events[$name][$i]);
			}
		}
	}

	public static function get($name)
	{
		return empty(self::$events[$name]) ? array(): self::$events[$name];
	}

	/*
	   Method: notify
	   If your event does not need to process the returned value from any
	   observers use this instead of get()
	*/
	public static function notify($name)
	{
		// removing event name from the arguments
		$args = func_num_args() > 1 ? array_slice(func_get_args(), 1): array();

		foreach (self::get($name) as $callback)
			call_user_func_array($callback, $args);
	}

} // end Event class
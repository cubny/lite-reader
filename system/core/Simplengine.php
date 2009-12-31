<?php
/**
 *  Class: SE
 *
 *  The Simple Engine main class launches the application and runs the system.
 *
 *  A route string can be a literal url such as '/pages/about' or contain
 *  wildcards (:any or :num) and/or regex like '/blog/:num' or '/page/:any'.
 *  The idea has taken from CodeIgniter.
 *
 *  > SE::addRoute(array(
 *  > '/' => 'page/index',
 *  > '/about' => 'page/about,
 *  > '/blog/:num' => 'blog/post/$1',
 *  > '/blog/:num/comment/:num/delete' => 'blog/deleteComment/$1/$2'
 *  > ));
 *
 *  Visiting /about/ would call PageController::about()
 *  visiting /blog/5 would call BlogController::post(5)
 *  visiting /blog/5/comment/42/delete would call BlogController::deleteComment(5,42)
 *
 *  Add your routes before dispatching because after this method it is too late.
*/

final class SE
{
	private static $routes = array();
	private static $params = array();
	private static $url = '';

	/**
	 *  Method: addRoute
	 *  Add a new route to the application
	 */
	public static function addRoute($route, $destination=null)
	{
		if ($destination != null)
			$route = array($route => $destination);

		self::$routes = array_merge(self::$routes, $route);
	}

	/*
	   Method: dispatch
	   Check for matching route then execute the action from the controller needed
	*/
	public static function dispatch($url=null)
	{
		Benchmark::start('dispatch');

		if ($url === null) $url = $_SERVER['QUERY_STRING'];

		// we populate the $_GET table
		if( $pos = strpos($url,'&') ) parse_str(substr($url, $pos), $_GET);

		// remove slashes (for route convention)
		$url = trim($url, '/');

		// removing the suffix for search engine static simulations
		if (URL_SUFFIX != null and ($pos_to_cut = strrpos($requested_url, $suffix)) !== false)
			$requested_url = substr($requested_url, 0, $pos_to_cut);

		self::$url = $url;

		if (empty($url))
		{
		    return self::executeAction(self::getController(), self::getAction(), self::getParams());
		}
		// do we even have any custom routing to deal with?
		else if (count(self::$routes) === 0)
		{
			self::$params = explode('/', $url);
			return self::executeAction(self::getController(), self::getAction(), self::getParams());
		}
		// is there a literal match? If so we're done
		else if (isset(self::$routes[$url]))
		{
			self::$params = explode('/', self::$routes[$url]);
			return self::executeAction(self::getController(), self::getAction(), self::getParams());
		}

		// loop through the route array looking for wildcards
		foreach (self::$routes as $rule => $route)
		{
			// convert wildcards to regex
			if (strpos($rule, ':') !== false)
				$rule = str_replace(':any', '(.+)', str_replace(':num', '([0-9]+)', $rule));

			// does the regex match?
			if (preg_match('#^'.$rule.'$#', $url))
			{
				// do we have a back-reference?
				if (strpos($route, '$') !== false and strpos($rule, '(') !== false)
					$route = preg_replace('#^'.$rule.'$#', $route, $url);

				self::$params = explode('/', $route);
				// we fund it, so we can break the loop now!
				return self::executeAction(self::getController(), self::getAction(), self::getParams());
			}
		}

		self::$params = explode('/', $url);
		return self::executeAction(self::getController(), self::getAction(), self::getParams());
	} // dispatch

	/*
	   Method: getCurrentUrl
	   Give the requested url
	*/
	public static function getCurrentUrl()
	{
		return self::$url;
	}

	/*
	   Method: getController
	   Give the controller used
	*/
	public static function getController()
	{
		return isset(self::$params[0]) ? self::$params[0]: DEFAULT_CONTROLLER;
	}

	/*
	   Method: getAction
	   Give the action performed
	*/
	public static function getAction()
	{
		return isset(self::$params[1]) ? self::$params[1]: DEFAULT_ACTION;
	}

	/*
	   Method: getParams
	   Give all additional parameters passed in the url
	*/
	public static function getParams()
	{
		return array_slice(self::$params, 2);
	}

	/*
	   Method: executeAction
	   Load the controller and execute the action
	*/
	public static function executeAction($controller, $action, $params)
	{
		$controller_class = Inflector::camelize($controller);
		$controller_class_name = $controller_class . 'Controller';

		// get a instance of that controller
		if (class_exists($controller_class_name))
			$controller = new $controller_class_name();

		if (!$controller instanceof Controller)
			throw new Exception("Class '{$controller_class_name}' does not extends Controller class!");

		Observer::notify('system.execute');
		Benchmark::start('execute');

		// execute the action
		$controller->execute($action, $params);
	}

} // end SE class

/*
   Function: page_not_found
   Display a 404 page not found and exit
*/
function page_not_found()
{
	Observer::notify('system.page_not_found');

	header("HTTP/1.0 404 Not Found");
	echo new View('404');

	Observer::notify('system.shutdown');
	exit;
}

// convert size in byte to the easiest humain readable size
function convert_size($num)
{
	if ($num >= 1073741824) $num = round($num / 1073741824 * 100) / 100 .' gb';
	else if ($num >= 1048576) $num = round($num / 1048576 * 100) / 100 .' mb';
	else if ($num >= 1024) $num = round($num / 1024 * 100) / 100 .' kb';
	else $num .= ' b';
	return $num;
}

// debug fonction displaying reversed trace ----------------------------------

function se_framework_exception_handler($e)
{
	if (!DEBUG) page_not_found();

	// display Profiler
	include SYSPATH.'/core/Profiler'.EXT;
	Profiler::displayTrace($e);
	Profiler::display();
}

set_exception_handler('se_framework_exception_handler');
Observer::notify('system.init');

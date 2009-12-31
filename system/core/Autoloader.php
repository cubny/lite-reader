<?php
/*
   Class: AutoLoader

   The AutoLoader class is an object oriented hook into PHP's __autoload functionality. You can add

   - Single Files AutoLoader::addFile('Blog','/path/to/Blog.php');
   - Multiple Files AutoLoader::addFile(array('Blog'=>'/path/to/Blog.php','Post'=>'/path/to/Post.php'));
   - Whole Folders AutoLoader::addFolder('path');

   When adding a whole folder each file should contain one class named the same as the file without ".php" (Blog => Blog.php)
*/
class AutoLoader
{
	protected static $files = array();
	protected static $folders = array();

	/*
	   Method: addFile
	   Add directly the file need for a specific class
	*/
	public static function addFile($class_name, $file=null)
	{
		if ($file == null)
			self::$files = array_merge(self::$files, $class_name);
		else
			self::$files[$class_name] = $file;
	}

	/*
	   Method: addFolder
	   Add a folder to soearch for the class
	*/
	public static function addFolder($folder)
	{
		if (!is_array($folder))
			$folder = array($folder);

		self::$folders = array_merge(self::$folders, $folder);
	}

	/**
	 *  Method: load
	 *  Will be called by PHP (trying to load the class)
	 */
	public static function load($class_name)
	{
		if (isset(self::$files[$class_name]))
		{
			require self::$files[$class_name];
			return;
		}
		else
		{
			foreach (self::$folders as $folder)
			{
				$file = $folder.$class_name.'.php';
				if (file_exists($file))
				{
					require $file;
					return;
				}
			}
		}
		throw new Exception("Class '{$class_name}' not found!");
	}

} // end AutoLoader class

if ( ! function_exists('__autoload')) {
	AutoLoader::addFolder(array(APPPATH.'/models/', APPPATH.'/controllers/'));

	function __autoload($class_name)
	{
		AutoLoader::load($class_name);
	}
}

/*
   Function: use_helper
   Load all classes and/or functions from the helper file(s)

   example:
   use_helper('I18n', 'Pagination');
*/
function use_helper()
{
	static $helpers = array();

	foreach (func_get_args() as $helper)
	{
		if (in_array($helper, $helpers)) continue;

		$helper_file = SYSPATH.'/helpers/'.DIRECTORY_SEPARATOR.$helper.EXT;

		if (!file_exists($helper_file))
			throw new Exception("Helper file '{$helper}' not found!");

		include $helper_file;
		$helpers[] = $helper;
	}
}

/*
   Function: use_model
   Load a model (faster then waiting for the__autoloader)

   example:
   use_model('Blog', 'Tag', 'User');
*/
function use_model()
{
	static $models = array();

	foreach (func_get_args() as $model)
	{
		if (in_array($model, $models)) continue;

		$model_file = APPPATH.'/models/'.$model.EXT;

		if (!file_exists($model_file))
			throw new Exception("Model file '{$model}' not found!");

		include $model_file;
		$models[] = $model;
	}
}
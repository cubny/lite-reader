<?php
/*
   Class: Inflector

   Static class used to pass from underscrore_string to CamelizeString to 
   underscored_one or to Humanize string
*/
final class Inflector 
{
	/*
	   Method: camelize
	   Format a string from an underscore string to an CamelizeSyntaxed.
	*/
	public static function camelize($string)
	{
		return str_replace(' ','',ucwords(str_replace('_',' ', $string)));
	}

	/*
	   Method: underscore
	   Format a string from an camelize string to an underscore_syntaxed.
	*/
	public static function underscore($string)
	{
		return strtolower(preg_replace('/(?<=\\w)([A-Z])/', '_\\1', $string));
	}

	/*
	   Method: humanize
	   Format a string from an underscore string to an Humanized syntaxed.
	*/
	public static function humanize($string)
	{
		return ucfirst(str_replace('_', ' ', $string));
	}

} // end Inflector class
<?php

class Folder
{
    static function getFiles($directory)
    {
	if(!is_dir($directory)) return false;
	
	$files = scandir($directory);

	foreach($files as $i => $value) {
		if (substr($value, 0, 1) == '.') {
			unset($files[$i]);
		}
	}
	return $files;
    }
}
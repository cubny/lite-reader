<?php

/**
 * Website application directory
 */
$application = 'app';

/**
 * Green system directory
 */
$system = 'system';

define('EXT', '.php');

// Define application and system paths
define('APPPATH',  dirname(__FILE__).'/'.$application.'/');
define('SYSPATH', dirname(__FILE__).'/'.$system.'/');
define('CACHEPATH', dirname(__FILE__).'/public/cache');

//define('DEFAULT_CONTROLLER', 'welcome');
//define('DEFAULT_ACTION', 'index');

//  Inititalize
require SYSPATH.'core/Bootstrap'.EXT;

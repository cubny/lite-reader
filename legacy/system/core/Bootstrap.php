<?php
require SYSPATH.'/helpers/Folder'.EXT;
/**
 * Bootstrap File
 */

/**
 * Getting all config files,
 * check application directory first then the system directory.
 */
$config_dir = 'config/';
$files = Folder::getFiles(SYSPATH.$config_dir);

foreach($files as $file){
    if(file_exists(APPPATH.$config_dir.$file)){
        include APPPATH.$config_dir.$file;
    }else{
        include SYSPATH.$config_dir.$file;
    }
}

/**
 * install db if not exists
 */
$db_file = VARPATH."/agg.db";
if(!is_file($db_file)){
  header('Location: '.$config['base_url'].'/install.php');
  exit;
}

// Set default controller and action
define('DEFAULT_CONTROLLER', $route['_default']);
define('DEFAULT_ACTION', 'index');

// Load core files
require SYSPATH.'/core/Model'.EXT;
require SYSPATH.'/core/View'.EXT;
require SYSPATH.'/core/Controller'.EXT;
require SYSPATH.'/helpers/Inflector'.EXT;
require SYSPATH.'/core/Benchmark'.EXT;
require SYSPATH.'/core/Observer'.EXT;
require SYSPATH.'/core/Autoloader'.EXT;
require SYSPATH.'/core/Simplengine'.EXT;

Benchmark::start('total');

// Set view suffix
define('VIEW_SUFFIX', $config['view_suffix']);

// Set base url
define('BASE_URL',  $config['base_url']);

// Set debugger
define('DEBUG', $config['debug']);
define('TABLE_PREFIX', $config['table_prefix']);

/**
 * Get all routes and process them
 */
if(isset($route)):
    foreach($route as $key => $value):
    SE::addRoute($key, $value);
    endforeach;
endif;
try{
  //$__CONN__=new PDO("mysql:host=localhost;dbname=aggregator","root","ugdvqh");
  $__CONN__=new PDO("sqlite:".VARPATH."agg.db");
}catch(PDOException $e){
  require SYSPATH.'/libraries/DoLite'.EXT;
  $__CONN__=new DoLite("sqlite:".VARPATH."agg.db");
}

/**
 * auto update database based on files in app/db/
 * naming conversions is app/db/update.[version].sql
 * version should be sequential
 *
 * fetch last version from db and execute the next version sql file in a loop
 *
 */
$result = $__CONN__->query("select value from config where key='version'");
if($result === false){
  $version = 0;
}else{
  $result = $result->fetch();
  $version = ((int) $result['value']);
}
$__CONN__->setAttribute(PDO::ATTR_EMULATE_PREPARES, 0);
while(file_exists(APPPATH.'db/update.'.($version+1).'.sql')){
  $version++;
  $query = file_get_contents(APPPATH.'db/update.'.$version.'.sql');
  $__CONN__->exec($query);
}
$__CONN__->setAttribute(PDO::ATTR_EMULATE_PREPARES, 1);
$__CONN__->exec("UPDATE config set value=".$version." where key='version'");

/* end of update */

Model::connection($__CONN__);
Model::getConnection()->exec("SET NAMES 'utf8'");
// Dispatch Simplengine
SE::dispatch();

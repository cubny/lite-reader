<?php
define('SYSPATH', dirname(__FILE__).'/system/');
define('VARPATH',  dirname(__FILE__).'/var/');
define('APPPATH',  dirname(__FILE__).'/app/');
define('CACHEPATH',  dirname(__FILE__).'/public/cache/');
require SYSPATH.'/helpers/Folder.php';
require SYSPATH.'/libraries/DoLite.php';

$config_dir = 'config/';
$files = Folder::getFiles(SYSPATH.$config_dir);

foreach($files as $file){
    if(file_exists(APPPATH.$config_dir.$file)){
        include APPPATH.$config_dir.$file;
    }else{
        include SYSPATH.$config_dir.$file;
    }
}

$messages = array();
if(!is_writable(VARPATH)){
  $messages[] = "please make ".VARPATH." writable for webserver";
}
if(!is_writable(CACHEPATH)){
  $messages[] = "please make ".CACHEPATH." writable for webserver";
}
if(count($messages)>0){
  echo implode("<br/>",$messages);
  exit;
}
if(!is_file(CACHEPATH.'images')){
  mkdir(CACHEPATH.'images');
}


try{
  $conn=new PDO("sqlite:".VARPATH."/agg.db");
}catch(PDOException $e){
  $conn=new DoLite("sqlite:".VARPATH."/agg.db");
}
$conn->setAttribute(PDO::ATTR_EMULATE_PREPARES, 0);
$query = file_get_contents(APPPATH.'db/install.sql');
$conn->exec($query);


header('Location: '.$config['base_url']);
exit;
?>

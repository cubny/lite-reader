<?php

$config['base_url'] = "http://".dirname($_SERVER['HTTP_HOST'].$_SERVER['SCRIPT_NAME'])."/";

$config['debug'] = true;

// Setting the suffix for the view pages located inside application directory views.
$config['view_suffix'] = '.phtml';
$config['table_prefix'] = '';

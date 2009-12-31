<?php

class IndexController extends Controller
{
  public function __construct(){
    use_model('Rss','Item');
  }
  public function index(){
     //var_dump($root);
    //$vars=array('title'=>(string)$root->channel->title);
    $rss=new Rss();
    $feeds=array();
    $data=$rss->getAllFeeds();
    //var_dump($rss->getColumns());
    //var_dump($rss->loadRss());
    $this->display("index",array('feeds'=>$data));
  }
  public function getItems($id){
    $item=new Item();
    $items=$item->getAllByRssId($id);
    header('Cache-Control: no-cache, must-revalidate'); 
    header('Expires: Mon, 26 Jul 1997 05:00:00 GMT');
    header('Content-type:text/json');
    echo json_encode($items);

  }
  public function getDesc($id){
    $item=new Item();
    $desc=$item->getDesc($id);
    header('Cache-Control: no-cache, must-revalidate'); 
    header('Expires: Mon, 26 Jul 1997 05:00:00 GMT');
    header('Content-type:text/json');
    echo json_encode($desc);

  }
  public function add($url){
    $rss=new Rss();
    //$url="http://news.google.com/news?pz=1&cf=all&ned=us&hl=en&topic=h&num=3&output=rss";
    //$url="http://www.irib.ir/rss/rssirib.xml";
    $data=$rss->add($url)->getAllFeeds();
    header('Cache-Control: no-cache, must-revalidate'); 
    header('Expires: Mon, 26 Jul 1997 05:00:00 GMT');
    header('Content-type:text/json');
    echo json_encode($data);
  }
}
  

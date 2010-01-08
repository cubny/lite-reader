<?php

class AggController extends Controller
{
  public function __construct(){
    use_model('Rss','Item');
  }

  private function _jsonify($data){
    header('Cache-Control: no-cache, must-revalidate'); 
    header('Expires: Mon, 26 Jul 1997 05:00:00 GMT');
    header('Content-type:text/json');
    return json_encode($data);
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
  public function getFav($id){
    $rss=new Rss();
    $link=$rss->load($id)->link;
  }
  public function getItems($id){
    $item=new Item();
    $items=$item->getAllByRssId($id);
    echo $this->_jsonify($items);
    
  }
  public function getDesc($id){
    $item=new Item();
    $desc=$item->getDesc($id);
    echo $this->_jsonify($desc);

  }
  public function add(){
    if(isset($_POST['url'])){
      $url=$_POST['url'];
    }

    $rss=new Rss();
    if($rss->add($url)!==false){
      $data=$rss->getLastFeed();
    }else{
      $feed=$rss->getByUrl($url);
      $data=array("error"=>"already Exists","feed"=>$feed->id);
    }
    echo $this->_jsonify($data);
  }
  public function feeds(){
    $rss=new Rss();
    echo $this->_jsonify($rss->getAllFeeds()); 
  }
  public function del($id){
    $rss=new Rss();
    $rss->id=$id;
    return $rss->delete();
  }
  public function update($id){
    $rss=new Rss();
    $rss->id=$id;
    $rss->updateItems();
    return $this->getItems($id);
  }
  public function fetch(){
    $url='http://feeds.simonwillison.net/swn-everything';
    $rss=new Rss();
    $data=$rss->fetch($url);
    var_dump($data);
  }
  public function test(){
    var_dump($_SERVER['HTTP_ACCEPT_ENCODING']);
    exit();
  }

}
  

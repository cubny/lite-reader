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
    $rss=new Rss();
    $feeds=array();
    $data=$rss->getAllFeeds();
    $this->display("index",array(
        'starred_count'=>Item::countFrom("Item","starred = 1"),
        'unread_count'=>Item::countFrom("Item","is_new = 1"),
        'feeds'=>$data
    ));
  }

  public function getFav($id){
    $rss=new Rss();
    $link=$rss->load($id)->link;
  }

  public function getItems($id){
    $item=new Item();
    switch($id){
      case "starred":
        $items=$item->getAllStarred();
        break;
      case "unread":
        $items=$item->getAllUnread();
        break;
      default:
        $items=$item->getAllByRssId($id);
    }
    $result = array();
    foreach($items as &$item){
      $doc = new DOMDocument();
      $doc->loadHTML($item->desc);
      $images = $doc->getElementsByTagName('img');
      foreach ($images as $tag) {
       $cache_name = 'agg/images/'.base64_encode($tag->getAttribute('src'));
       $tag->setAttribute('src',"public/images/grey.gif");
       $tag->setAttribute('data-original',$cache_name);
       $tag->setAttribute('class','lazy not-loaded');
      }
      $item->desc = $doc->saveHTML();
    }
    echo $this->_jsonify($items);
  }

  public function images($url_encode){
    $image_path = CACHEPATH."/images/";
    if(!is_file($image_path.$url_encode)){
      $image = file_get_contents(base64_decode($url_encode));
      file_put_contents($image_path.$url_encode,$image);
    }
    $filename = basename(base64_decode($url_encode));
    $size = getimagesize($image_path.$url_encode);
    header('Content-type: ' . $size['mime']);
    echo  file_get_contents($image_path.$url_encode);
  }

  public function mark_read_all($rss_id){
    $item=new Item();
    $result=$item->make_read_all($rss_id);
    echo $result;
  }

  public function mark_unread_all($rss_id){
    $item=new Item();
    $result=$item->make_unread_all($rss_id);
    echo $result;
  }

  public function make_read($id){
    $item=new Item();
    $result=$item->make_read($id);
    echo $result;
  }

  public function make_unread($id){
    $item=new Item();
    $result=$item->make_unread($id);
    echo $result;
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
    $result = $rss->delete();
    die(var_dump($result));
  }
  public function update_all(){
    $rss = new Rss();
    $data = array();
    $feeds=$rss->getAllFeeds();
    foreach($feeds as $feed){
      $rss->id = $feed->id;
      $rss->updateItems();
      $item=new Item();
      $items=$item->getAllByRssId($feed->id);
      $data[$feed->id] = $items;
    }
    echo $this->_jsonify($data);
  }
  public function update($id){
    $rss=new Rss();
    $rss->id=$id;
    $rss->updateItems();
    return $this->getItems($id);
  }
  public function star($item_id){
    $item=new Item();
    $result=$item->make_starred($item_id);
    echo $result;
  }
  public function unstar($item_id){
    $item=new Item();
    $result=$item->make_unstarred($item_id);
    echo $result;
  }
}
  

<?php
//include SYSPATH."/libraries/rayfeedreader.php";
require SYSPATH."/libraries/SimplePie.compiled.php";
class Rss extends Model
{
  const TABLE_NAME='rss';

  public $url;
  public $title;
  public $desc;
  public $link;
  public $lang;
  public $updated_at;

  public function __construct(){
    use_helper('I18n');
  }
  public function fetchFeed($url){
    $feed = new SimplePie();
    $feed->set_feed_url($url);
    $feed->init();
    $feed->handle_content_type();

    return $feed;
  }
  public function updateItems(){
   $feed=self::findByIdFrom(get_class($this),$this->id);
   $url=$feed->url;
   $data=$this->fetchFeed($url);
   $items=$data->get_items();

   foreach($items as $item){
     $itemModel=new Item();
     if(!$itemModel->exists($item->get_link())){
      $this->_addItem($item,$itemModel);
     }
     unset($itemModel);
   }
  }
  private function _addItem($item,$itemModel){
    $itemModel->title=$item->get_title();
    $itemModel->link=$item->get_link();
    $itemModel->rss_id=$this->id;
    $itemModel->is_new=1;
    $itemModel->desc = $item->get_content();
    $itemModel->dir = $itemModel->desc;
    $itemModel->save();

    /*$descs = array('content:encoded','content','description','title');
    foreach($descs as $key){
      if(!empty($item[$key])){
        $itemModel->desc = $item[$key];
        $itemModel->save();
        return;
      }
    }*/
  }
  public function add($url){
    //$url=urldecode($url);
    if($this->urlExists($url)>0)return false;
    $data=$this->fetchFeed($url);
    $this->url=$url;
    $this->title=$data->get_title();
    $this->desc=$data->get_description();
    $this->link=$data->get_link();
    $this->lang=is_null($data->get_language())?$data->get_language():'en';
    $this->updated_at=(string) date("Y-m-d H:i");
    $this->save();
    $items=$data->get_items();
    foreach($items as $item){
      $itemModel=new Item();
      $this->_addItem($item,$itemModel);
      unset($itemModel);
    }
    return $this;
  }
  public function getLastFeed(){
    return $this->query('SELECT rss.id as id,rss.title as title, rss.link as link, rss.url as url, (SELECT count(item.is_new) from '.Item::TABLE_NAME.' WHERE item.rss_id=rss.id AND item.is_new="1") as unread from '.self::TABLE_NAME.' as rss ORDER BY rss.id DESC limit 1')->fetchAll(self::FETCH_OBJ);
  }
  public function getAllFeeds(){
    return $this->query('SELECT rss.id as id,rss.title as title, rss.link as link, rss.url as url, (SELECT count(item.is_new) from '.Item::TABLE_NAME.' WHERE item.rss_id=rss.id AND item.is_new="1") as unread from '.self::TABLE_NAME.' as rss ORDER BY rss.id')->fetchAll(self::FETCH_OBJ);
  }
  public function urlExists($url){
    return self::countFrom(get_class($this),"url = ?",array($url));
  }
  public function getByUrl($url){
    //$url=urldecode($url);
    return self::findOneFrom(get_class($this),"url = ?",array($url));
  }

  public function beforeDelete(){
    $item=new Item();
    return $item->deleteByRssId($this->id);
  }
 
}

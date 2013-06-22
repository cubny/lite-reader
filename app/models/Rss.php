<?php
include SYSPATH."/libraries/rayfeedreader.php";
class Rss extends Model
{
  const TABLE_NAME='rss';

  public $url;
  public $title;
  public $desc;
  public $link;
  public $lang;
  public $updated_at;

  public function fetchFeed($url){
    $ray=new RayFeedReader(array(
      'url'=>$url,
      'httpClient'=>'php',
    ));
    $data=$ray->parse()->getData();
    if(!isset($data['items'])){
        $data['items'] = array();
    }
    $data['items']=array_reverse($data['items']);
    return $data;
  }
  public function updateItems(){
   $feed=self::findByIdFrom(get_class($this),$this->id);
   $url=$feed->url;
   $data=$this->fetchFeed($url);
   $items=$data['items'];
   foreach($items as $item){
     $itemModel=new Item();
     if(!$itemModel->exists($item['link'])){
      $this->_addItem($item,$itemModel);
     }
     unset($itemModel);
   }
  }
  private function _addItem($item,$itemModel){
    $itemModel->title=$item['title'];
    $itemModel->link=$item['link'];
    $itemModel->rss_id=$this->id;
    $itemModel->is_new=1;
    if(empty($item['description'])){
      if(empty($item['content'])){
        $itemModel->desc=$item['title'];
      }else{
        $itemModel->desc = $item['content'];
      }
    }else{
      $itemModel->desc=$item['description'];
    }
    $itemModel->save();
  }
  public function add($url){
    //$url=urldecode($url);
    if($this->urlExists($url)>0)return false;
    $data=$this->fetchFeed($url);
    $this->url=$url;
    $this->title=$data['title'];
    $this->desc=$data['description'];
    $this->link=$data['link'];
    $this->lang=isset($data['language'])?$data['language']:'en';
    $this->updated_at=(string) date("Y-m-d H:i");
    $this->save();
    $items=$data['items'];
    foreach($items as $item){
      $itemModel=new Item();
      $this->_addItem($item,$itemModel);
      unset($itemModel);
    }
    return $this;
  }
  public function getLastFeed(){
    return $this->query('SELECT rss.id as id,rss.title as title, rss.link as link, (SELECT count(item.is_new) from '.Item::TABLE_NAME.' WHERE item.rss_id=rss.id AND item.is_new="1") as unread from '.self::TABLE_NAME.' as rss ORDER BY rss.id DESC limit 1')->fetchAll(self::FETCH_OBJ);
  }
  public function getAllFeeds(){
    return $this->query('SELECT rss.id as id,rss.title as title, rss.link as link, (SELECT count(item.is_new) from '.Item::TABLE_NAME.' WHERE item.rss_id=rss.id AND item.is_new="1") as unread from '.self::TABLE_NAME.' as rss')->fetchAll(self::FETCH_OBJ);
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

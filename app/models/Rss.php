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

  public function add($url){
    $url=urldecode($url);
    if($this->urlExists($url)>0)return;
    $ray=new RayFeedReader(array(
      'url'=>$url,
      'httpClient'=>'SimpleXML',
      ));
    $data=$ray->parse()->getData();
    $this->url=$url;
    $this->title=$data['title'];
    $this->desc=$data['description'];
    $this->link=$data['link'];
    $this->lang=isset($data['language'])?$data['language']:'en';
    $thiis->updated_at=date("Y-m-d H:i");
    $this->save();
    $items=$data['items'];
    foreach($items as $item){
      $itemModel=new Item();
      $itemModel->title=$item['title'];
      $itemModel->link=$item['link'];
      $itemModel->rss_id=$this->id;
      $itemModel->is_new=1;
      if(empty($item['description'])){
        $itemModel->desc=$item['title'];
      }else{
        $itemModel->desc=$item['description'];
      }
      $itemModel->save();
      unset($itemModel);
    }
    return $this;
  }
  public function getLastFeed(){
    return $this->query('SELECT rss.id as id,rss.title as title,(SELECT count(item.is_new) from '.Item::TABLE_NAME.' WHERE item.rss_id=rss.id AND item.is_new="1") as unread from '.self::TABLE_NAME.' as rss ORDER BY rss.id DESC limit 1')->fetchAll(self::FETCH_OBJ);
  }
  public function getAllFeeds(){
    return $this->query('SELECT rss.id as id,rss.title as title,(SELECT count(item.is_new) from '.Item::TABLE_NAME.' WHERE item.rss_id=rss.id AND item.is_new="1") as unread from '.self::TABLE_NAME.' as rss')->fetchAll(self::FETCH_OBJ);
  }
  public function urlExists($url){
    return self::countFrom(get_class($this),"url = ?",array($url));
  }

  public function beforeDelete(){
    $item=new Item();
    return $item->deleteByRssId($this->id);
  }
 
}

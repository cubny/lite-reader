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
    $this->save();
    $items=$data['items'];
    foreach($items as $item){
      $itemModel=new Item();
      $itemModel->title=$item['title'];
      $itemModel->link=$item['link'];
      $itemModel->rss_id=$this->id;
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
  public function getAllFeeds(){
    return $this->query('SELECT id,title from '.self::TABLE_NAME)->fetchAll(self::FETCH_OBJ);
  }
  public function urlExists($url){
    return self::countFrom(get_class($this),"url = ?",array($url));
  }
  
}

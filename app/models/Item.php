<?php
class Item extends Model
{
  const TABLE_NAME='item';
  public $title;
  public $link;
  public $rss_id;
  public $desc;
  public $is_new;

  public function getAllByRssId($id){
    return $this->query("SELECT * from ".self::TABLE_NAME." where rss_id=$id")->fetchAll(self::FETCH_OBJ);
  }
  public function getDesc($id){
    $this->id=$id;
    self::update(self::TABLE_NAME,array("is_new"=>"0"),"id=?",array($id));
    return $this->query("SELECT title,link,desc,is_new from ".self::TABLE_NAME." where id=$id")->fetchObject();
  }
  public function deleteByRssId($rss_id){
    return self::deleteWhere(get_class($this),"rss_id = ?",array($rss_id));
  }

}


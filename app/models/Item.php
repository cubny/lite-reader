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
    return $this->query("SELECT * from ".self::TABLE_NAME." where rss_id=$id ORDER BY ID")->fetchAll(self::FETCH_OBJ);
  }
  public function getAllStarred(){
    return $this->query("SELECT * from ".self::TABLE_NAME." where starred=1 ORDER BY ID")->fetchAll(self::FETCH_OBJ);
  }
  public function getAllUnread(){
    return $this->query("SELECT * from ".self::TABLE_NAME." where is_new=1 ORDER BY ID")->fetchAll(self::FETCH_OBJ);
  }
  public function make_unread_all($rss_id){
    return self::update(self::TABLE_NAME,array("is_new"=>"1"),"rss_id=?",array($rss_id));
  }
  public function make_read_all($rss_id){
    return self::update(self::TABLE_NAME,array("is_new"=>"0"),"rss_id=?",array($rss_id));
  }
  public function make_read($id){
    $this->id=$id;
    return self::update(self::TABLE_NAME,array("is_new"=>"0"),"id=?",array($id));
  }
  public function make_unread($id){
    $this->id=$id;
    return self::update(self::TABLE_NAME,array("is_new"=>"1"),"id=?",array($id));
  }
  public function make_starred($id){
    $this->id=$id;
    return self::update(self::TABLE_NAME,array("starred"=>1),"id=?",array($id));
  }
  public function make_unstarred($id){
    $this->id=$id;
    return self::update(self::TABLE_NAME,array("starred"=>0),"id=?",array($id));
  }
  public function getDesc($id){
    $this->id=$id;
    self::update(self::TABLE_NAME,array("is_new"=>"0"),"id=?",array($id));
    return $this->query("SELECT title,link,desc,is_new from ".self::TABLE_NAME." where id=$id")->fetchObject();
  }
  public function deleteByRssId($rss_id){
    return self::deleteWhere(get_class($this),"rss_id = ?",array($rss_id));
  }
  public function exists($link){
    return self::countFrom(get_class($this),"link = ?",array($link));
  }

}


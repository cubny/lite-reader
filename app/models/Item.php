<?php
class Item extends Model
{
  const TABLE_NAME='item';
  public $title;
  public $link;
  public $rss_id;
  public $desc;

  public function getAllByRssId($id){
    return $this->query("SELECT * from ".self::TABLE_NAME." where rss_id=$id")->fetchAll(self::FETCH_OBJ);
  }
  public function getDesc($id){
    $this->id=$id;
    return $this->query("SELECT title,link,desc from ".self::TABLE_NAME." where id=$id")->fetchObject();
  }
}


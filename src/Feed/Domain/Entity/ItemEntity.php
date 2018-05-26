<?php

namespace LiteReader\Feed\Domain\Entity;

use LiteReader\Feed\Domain\ValueObjects\{
    Id,
    Title,
    Description,
    Url,
    Unread,
    PublishedAt
};

final class ItemEntity
{
    /**
     * id 
     *
     * @var Id
     * @access private
     */
    private $id;

    /**
     * title 
     *
     * @var Title 
     * @access private
     */
    private $title; 

    /**
     * discription 
     *
     * @var Description 
     * @access private
     */
    private $discription;

    /**
     * url 
     *
     * @var Url 
     * @access private
     */
    private $url;

    /**
     * unread 
     *
     * @var Unread 
     * @access private
     */
    private $unread;

    /**
     * publishedAt 
     *
     * @var PublishedAt 
     * @access private
     */
    private $publishedAt;

    public function __construct(
        Id $id,
        Title $title,
        Description $description,
        Url $url,
        Unread $unread,
        PublishedAt $publishedAt
    ) {
        $this->id = $id;
        $this->title = $title;
        $this->description = $description;
        $this->url = $url;
        $this->unread = $unread;
        $this->publishedAt = $publishedAt;
    }
}

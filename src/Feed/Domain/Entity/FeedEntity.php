<?php

namespace LiteReader\Feed\Domain\Entity;

use LiteReader\Feed\Domain\ValueObjects\{
    Id,
    Title,
    Description,
    Url,
    UpdatedAt
};

final class FeedEntity
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
     * updatedAt 
     *
     * @var UpdatedAt 
     * @access private
     */
    private $updatedAt;

    private $items;

    public function __construct(
        Id $id,
        Title $title,
        Description $description,
        Url $url
    ) {
        $this->id = $id;
        $this->title = $title;
        $this->description = $description;
        $this->url = $url;
    }

    public function addNewItems()
    {
        $this->updatedAt = new UpdatedAt();
    }
}

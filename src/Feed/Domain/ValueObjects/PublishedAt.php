<?php

namespace LiteReader\Feed\Domain\ValueObjects;

final class PublishedAt
{
    /**
     *
     * @var DateTime 
     */
    private $updateAt;

    /**
     * Constructoror
     *
     * @param DateTime $publishedAt
     * @access public
     * @return void
     */
    public function __constructor(DateTime $publishedAt)
    {
        $this->publishedAt = $publishedAt;
    }

    /**
     * getValue
     *
     * @access public
     * @return DateTime
     */
    public function getValue(): DateTime
    {
        return $this->publishedAt;
    }
}

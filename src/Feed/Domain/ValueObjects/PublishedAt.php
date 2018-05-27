<?php

namespace LiteReader\Feed\Domain\ValueObjects;

use \Datetime;

final class PublishedAt
{
    /**
     *
     * @var DateTime
     */
    private $publishedAt;

    /**
     * Constructoror
     *
     * @param DateTime $publishedAt
     * @access public
     * @return void
     */
    public function __construct(?DateTime $publishedAt = null)
    {
        if ($publishedAt === null) {
            $publishedAt = new DateTime();
        }
        $this->publishedAt = $publishedAt;
    }

    /**
     * getValue
     *
     * @access public
     * @return \DateTime
     */
    public function getValue(): \DateTime
    {
        return $this->publishedAt;
    }
}

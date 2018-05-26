<?php

namespace LiteReader\Feed\Domain\ValueObjects;

final class UpdatedAt
{
    /**
     *
     * @var DateTime 
     */
    private $updateAt;

    /**
     * Constructoror
     *
     * @param DateTime $updatedAt
     * @access public
     * @return void
     */
    public function __constructor(DateTime $updatedAt)
    {
        $this->updatedAt = $updatedAt;
    }

    /**
     * getValue
     *
     * @access public
     * @return DateTime
     */
    public function getValue(): DateTime
    {
        return $this->updatedAt;
    }
}

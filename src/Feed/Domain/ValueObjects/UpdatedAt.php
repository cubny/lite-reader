<?php

namespace LiteReader\Feed\Domain\ValueObjects;

use DateTime;

final class UpdatedAt
{
    /**
     *
     * @var DateTime
     */
    private $updatedAt;

    /**
     * Constructoror
     *
     * @param DateTime $updatedAt
     * @access public
     * @return void
     */
    public function __construct(?DateTime $updatedAt = null)
    {
        if ($updatedAt === null) {
            $updatedAt = new DateTime();
        }
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

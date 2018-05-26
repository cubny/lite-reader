<?php

namespace LiteReader\Feed\Domain\ValueObjects;

final class Unread
{
    /**
     * unread 
     *
     * @var bool
     * @access private
     */
    private $unread;

    /**
     * Constructor
     *
     * @param int $unread
     * @access public
     */
    public function __construct(bool $unread)
    {
        $this->unread = $unread;
    }

    /**
     * getValue
     *
     * @access public
     * @return bool
     */
    public function getValue(): bool
    {
        return $this->unread;
    }
}

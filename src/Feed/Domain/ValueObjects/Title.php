<?php

namespace LiteReader\Feed\Domain\ValueObjects;

final class Title
{
    /**
     *
     * @var string 
     */
    private $title;

    /**
     * Constructoror
     *
     * @param string $title
     * @access public
     * @return void
     */
    public function __construct(string $title)
    {
        $this->title = $title;
    }

    /**
     * getValue
     *
     * @access public
     * @return string
     */
    public function getValue(): string
    {
        return $this->title;
    }
}

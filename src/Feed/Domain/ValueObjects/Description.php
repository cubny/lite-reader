<?php

namespace LiteReader\Feed\Domain\ValueObjects;

final class Description
{
    /**
     *
     * @var string 
     */
    private $description;

    /**
     * Constructoror
     *
     * @param string $description
     * @access public
     * @return void
     */
    public function __construct(string $description)
    {
        $this->description = $description;
    }

    /**
     * getValue
     *
     * @access public
     * @return string
     */
    public function getValue(): string
    {
        return $this->description;
    }
}

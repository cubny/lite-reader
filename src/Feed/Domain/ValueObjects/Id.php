<?php

namespace LiteReader\Feed\Domain\ValueObjects;

final class Id
{
    /**
     * id 
     *
     * @var int 
     * @access private
     */
    private $id;

    /**
     * Constructor
     *
     * @param int $id
     * @access public
     */
    public function __construct(int $id)
    {
        $this->id = $id;
    }

    /**
     * getValue
     *
     * @access public
     * @return int
     */
    public function getValue(): int
    {
        return $this->id;
    }
}

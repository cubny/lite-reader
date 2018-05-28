<?php

namespace LiteReader\Feed\Domain\Collection;

use \InvalidArgumentException;
use LiteReader\Common\Domain\Collection\ArrayCollection;
use LiteReader\Feed\Domain\Entity\ItemEntity;

class ItemArrayCollection extends ArrayCollection
{
    public function __construct(array $items=[])
    {
        if ($this->countNotItemEntity($items) > 0 ) {
            throw new InvalidArgumentException("All Elements Must Be ItemEntity");
        }
        parent::__construct($items);
    }

    private function countNotItemEntity(array $items): int
    {
        return count(array_filter($items, [ $this, 'isNotItemEntity'] ));
    }

    private function isNotItemEntity($element): bool
    {
        return !($element instanceof ItemEntity);
    }
}

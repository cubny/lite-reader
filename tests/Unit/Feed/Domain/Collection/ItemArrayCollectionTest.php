<?php
/**
 * Created by PhpStorm.
 * User: cubny
 * Date: 5/28/2018 AD
 * Time: 10:40
 */

namespace Tests\Unit\Feed\Domain\Collection;

use InvalidArgumentException;
use LiteReader\Common\Domain\Collection\ArrayCollection;
use LiteReader\Feed\Domain\Collection\ItemArrayCollection;
use PHPUnit\Framework\TestCase;

class ItemArrayCollectionTest extends TestCase
{
    public function testRejectsWhenItemEntityIsNotGiven()
    {
        $this->expectException(InvalidArgumentException::class);
        new ItemArrayCollection([1,2,3]);
    }

    public function testRetunrsArrayCollectionWhenCalledWithNoParams()
    {
        $items = new ItemArrayCollection();
        $this->assertInstanceOf(ArrayCollection::class, $items);
    }
}

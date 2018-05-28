<?php 

namespace LiteReader\Unit\Feed\Domain\Entity;

use \InvalidArgumentException;
use \TypeError;
use LiteReader\Feed\Domain\Collection\ItemArrayCollection;
use LiteReader\Feed\Domain\Entity\FeedEntity;
use LiteReader\Feed\Domain\ValueObjects\{
    Id,
    Title,
    Description,
    Url,
    UpdatedAt
};
use PHPUnit\Framework\TestCase;

final class FeedEntityTest extends TestCase
{
    public function setUp()
    {
        $this->feed = new FeedEntity(
            new Id(1),
            new Title('title'),
            new Description('description description'),
            new Url('http://newfeed.com/rss')
        );
    }

    /**
     * @dataProvider attributeValueObjectItems
     * @param string $type
     * @param string $attribute
     */
    public function testAttributeIsInstanceOfValueObjects(string $type, string $attribute)
    {
        $this->assertAttributeInstanceOf($type, $attribute, $this->feed);
    }

    public function attributeValueObjectItems()
    {
        return [
            "id" => [Id::class, 'id'],
            "title" => [Title::class, 'title'],
            "description" => [Description::class, 'description'],
            "Url" => [Url::class, 'url']
        ];
    }

    public function testAddItemsRejectsArray()
    {
        $this->expectException(TypeError::class);
        $this->feed->addItems([1,2,3]);
    }

    public function testAddItemsAcceptsItemArrayCollection()
    {
        $itemsCollection = $this->prophesize(ItemArrayCollection::class);
        $this->feed->addItems($itemsCollection->reveal());
        $this->assertAttributeInstanceOf(ItemArrayCollection::class, 'items', $this->feed);
    }

}

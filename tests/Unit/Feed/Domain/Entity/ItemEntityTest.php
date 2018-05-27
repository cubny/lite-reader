<?php 

namespace LiteReader\Unit\Feed\Domain\Entity;

use LiteReader\Feed\Domain\Entity\ItemEntity;
use LiteReader\Feed\Domain\ValueObjects\{
    Id,
    Title,
    Description,
    Url,
    Unread,
    PublishedAt
};
use PHPUnit\Framework\TestCase;

final class ItemEntityTest extends TestCase
{
    private $item;

    public function setUp()
    {
        $this->item = new ItemEntity(
            new Id(1),
            new Title('title'),
            new Description('description description'),
            new Url('http://newfeed.com/rss'),
            new Unread(true),
            new PublishedAt(new \DateTime('now'))
        );
    }

    /**
     * @dataProvider attributeValueObjectItems
     * @param string $type
     * @param string $attribute
     */
    public function testAttributeIsInstanceOfValueObjects(string $type, string $attribute)
    {
        $this->assertAttributeInstanceOf($type, $attribute, $this->item);
    }

    public function attributeValueObjectItems()
    {
        return [
            "id" => [Id::class, 'id'],
            "title" => [Title::class, 'title'],
            "description" => [Description::class, 'description'],
            "Url" => [Url::class, 'url'],
            "unread" => [Unread::class, 'unread'],
            "publishedAt" => [PublishedAt::class, 'publishedAt']
        ];
    }
}

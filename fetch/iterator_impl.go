package fetch

// AUTO GENERATED
// see iterator_impl.gen.py for details

import (
	"errors"
	"fmt"
	"github.com/go-fed/activity/streams/vocab"
	"net/url"
)

func begin(iterable interface{}) (Iter, error) {
	if iterable == nil {
		return nil, errors.New("nil argument")
	}

	switch v := iterable.(type) {
	case vocab.ActivityStreamsActorProperty:
		return iter_ActivityStreamsActorPropertyIterator{
			p:  v,
			it: v.Begin(),
		}, nil

	case vocab.ActivityStreamsAnyOfProperty:
		return iter_ActivityStreamsAnyOfPropertyIterator{
			p:  v,
			it: v.Begin(),
		}, nil

	case vocab.ActivityStreamsAttachmentProperty:
		return iter_ActivityStreamsAttachmentPropertyIterator{
			p:  v,
			it: v.Begin(),
		}, nil

	case vocab.ActivityStreamsAttributedToProperty:
		return iter_ActivityStreamsAttributedToPropertyIterator{
			p:  v,
			it: v.Begin(),
		}, nil

	case vocab.ActivityStreamsAudienceProperty:
		return iter_ActivityStreamsAudiencePropertyIterator{
			p:  v,
			it: v.Begin(),
		}, nil

	case vocab.ActivityStreamsBccProperty:
		return iter_ActivityStreamsBccPropertyIterator{
			p:  v,
			it: v.Begin(),
		}, nil

	case vocab.ActivityStreamsBtoProperty:
		return iter_ActivityStreamsBtoPropertyIterator{
			p:  v,
			it: v.Begin(),
		}, nil

	case vocab.ActivityStreamsCcProperty:
		return iter_ActivityStreamsCcPropertyIterator{
			p:  v,
			it: v.Begin(),
		}, nil

	case vocab.ActivityStreamsClosedProperty:
		return iter_ActivityStreamsClosedPropertyIterator{
			p:  v,
			it: v.Begin(),
		}, nil

	case vocab.ActivityStreamsContextProperty:
		return iter_ActivityStreamsContextPropertyIterator{
			p:  v,
			it: v.Begin(),
		}, nil

	case vocab.ActivityStreamsFormerTypeProperty:
		return iter_ActivityStreamsFormerTypePropertyIterator{
			p:  v,
			it: v.Begin(),
		}, nil

	case vocab.ActivityStreamsGeneratorProperty:
		return iter_ActivityStreamsGeneratorPropertyIterator{
			p:  v,
			it: v.Begin(),
		}, nil

	case vocab.ActivityStreamsIconProperty:
		return iter_ActivityStreamsIconPropertyIterator{
			p:  v,
			it: v.Begin(),
		}, nil

	case vocab.ActivityStreamsImageProperty:
		return iter_ActivityStreamsImagePropertyIterator{
			p:  v,
			it: v.Begin(),
		}, nil

	case vocab.ActivityStreamsInReplyToProperty:
		return iter_ActivityStreamsInReplyToPropertyIterator{
			p:  v,
			it: v.Begin(),
		}, nil

	case vocab.ActivityStreamsInstrumentProperty:
		return iter_ActivityStreamsInstrumentPropertyIterator{
			p:  v,
			it: v.Begin(),
		}, nil

	case vocab.ActivityStreamsItemsProperty:
		return iter_ActivityStreamsItemsPropertyIterator{
			p:  v,
			it: v.Begin(),
		}, nil

	case vocab.ActivityStreamsLocationProperty:
		return iter_ActivityStreamsLocationPropertyIterator{
			p:  v,
			it: v.Begin(),
		}, nil

	case vocab.ActivityStreamsObjectProperty:
		return iter_ActivityStreamsObjectPropertyIterator{
			p:  v,
			it: v.Begin(),
		}, nil

	case vocab.ActivityStreamsOneOfProperty:
		return iter_ActivityStreamsOneOfPropertyIterator{
			p:  v,
			it: v.Begin(),
		}, nil

	case vocab.ActivityStreamsOrderedItemsProperty:
		return iter_ActivityStreamsOrderedItemsPropertyIterator{
			p:  v,
			it: v.Begin(),
		}, nil

	case vocab.ActivityStreamsOriginProperty:
		return iter_ActivityStreamsOriginPropertyIterator{
			p:  v,
			it: v.Begin(),
		}, nil

	case vocab.ActivityStreamsPreviewProperty:
		return iter_ActivityStreamsPreviewPropertyIterator{
			p:  v,
			it: v.Begin(),
		}, nil

	case vocab.ActivityStreamsRelationshipProperty:
		return iter_ActivityStreamsRelationshipPropertyIterator{
			p:  v,
			it: v.Begin(),
		}, nil

	case vocab.ActivityStreamsResultProperty:
		return iter_ActivityStreamsResultPropertyIterator{
			p:  v,
			it: v.Begin(),
		}, nil

	case vocab.ActivityStreamsStreamsProperty:
		return iter_ActivityStreamsStreamsPropertyIterator{
			p:  v,
			it: v.Begin(),
		}, nil

	case vocab.ActivityStreamsTagProperty:
		return iter_ActivityStreamsTagPropertyIterator{
			p:  v,
			it: v.Begin(),
		}, nil

	case vocab.ActivityStreamsTargetProperty:
		return iter_ActivityStreamsTargetPropertyIterator{
			p:  v,
			it: v.Begin(),
		}, nil

	case vocab.ActivityStreamsToProperty:
		return iter_ActivityStreamsToPropertyIterator{
			p:  v,
			it: v.Begin(),
		}, nil

	case vocab.ActivityStreamsUrlProperty:
		return iter_ActivityStreamsUrlPropertyIterator{
			p:  v,
			it: v.Begin(),
		}, nil

	case vocab.W3IDSecurityV1PublicKeyProperty:
		return iter_W3IDSecurityV1PublicKeyPropertyIterator{
			p:  v,
			it: v.Begin(),
		}, nil

	default:
		return nil, fmt.Errorf("type=%T not supported", iterable)
	}
}

type iter_ActivityStreamsActorPropertyIterator struct {
	p  vocab.ActivityStreamsActorProperty
	it vocab.ActivityStreamsActorPropertyIterator
}

func (i iter_ActivityStreamsActorPropertyIterator) HasAny() bool {
	return i.it.HasAny()
}

func (i iter_ActivityStreamsActorPropertyIterator) IsIRI() bool {
	return i.it.IsIRI()
}

func (i iter_ActivityStreamsActorPropertyIterator) GetIRI() *url.URL {
	return i.it.GetIRI()
}

func (i iter_ActivityStreamsActorPropertyIterator) GetType() vocab.Type {
	return i.it.GetType()
}

func (i iter_ActivityStreamsActorPropertyIterator) Next() Iter {
	return iter_ActivityStreamsActorPropertyIterator{
		p:  i.p,
		it: i.it.Next(),
	}
}

func (i iter_ActivityStreamsActorPropertyIterator) End() Iter {
	return iter_ActivityStreamsActorPropertyIterator{
		p:  i.p,
		it: i.p.End(),
	}
}

type iter_ActivityStreamsAnyOfPropertyIterator struct {
	p  vocab.ActivityStreamsAnyOfProperty
	it vocab.ActivityStreamsAnyOfPropertyIterator
}

func (i iter_ActivityStreamsAnyOfPropertyIterator) HasAny() bool {
	return i.it.HasAny()
}

func (i iter_ActivityStreamsAnyOfPropertyIterator) IsIRI() bool {
	return i.it.IsIRI()
}

func (i iter_ActivityStreamsAnyOfPropertyIterator) GetIRI() *url.URL {
	return i.it.GetIRI()
}

func (i iter_ActivityStreamsAnyOfPropertyIterator) GetType() vocab.Type {
	return i.it.GetType()
}

func (i iter_ActivityStreamsAnyOfPropertyIterator) Next() Iter {
	return iter_ActivityStreamsAnyOfPropertyIterator{
		p:  i.p,
		it: i.it.Next(),
	}
}

func (i iter_ActivityStreamsAnyOfPropertyIterator) End() Iter {
	return iter_ActivityStreamsAnyOfPropertyIterator{
		p:  i.p,
		it: i.p.End(),
	}
}

type iter_ActivityStreamsAttachmentPropertyIterator struct {
	p  vocab.ActivityStreamsAttachmentProperty
	it vocab.ActivityStreamsAttachmentPropertyIterator
}

func (i iter_ActivityStreamsAttachmentPropertyIterator) HasAny() bool {
	return i.it.HasAny()
}

func (i iter_ActivityStreamsAttachmentPropertyIterator) IsIRI() bool {
	return i.it.IsIRI()
}

func (i iter_ActivityStreamsAttachmentPropertyIterator) GetIRI() *url.URL {
	return i.it.GetIRI()
}

func (i iter_ActivityStreamsAttachmentPropertyIterator) GetType() vocab.Type {
	return i.it.GetType()
}

func (i iter_ActivityStreamsAttachmentPropertyIterator) Next() Iter {
	return iter_ActivityStreamsAttachmentPropertyIterator{
		p:  i.p,
		it: i.it.Next(),
	}
}

func (i iter_ActivityStreamsAttachmentPropertyIterator) End() Iter {
	return iter_ActivityStreamsAttachmentPropertyIterator{
		p:  i.p,
		it: i.p.End(),
	}
}

type iter_ActivityStreamsAttributedToPropertyIterator struct {
	p  vocab.ActivityStreamsAttributedToProperty
	it vocab.ActivityStreamsAttributedToPropertyIterator
}

func (i iter_ActivityStreamsAttributedToPropertyIterator) HasAny() bool {
	return i.it.HasAny()
}

func (i iter_ActivityStreamsAttributedToPropertyIterator) IsIRI() bool {
	return i.it.IsIRI()
}

func (i iter_ActivityStreamsAttributedToPropertyIterator) GetIRI() *url.URL {
	return i.it.GetIRI()
}

func (i iter_ActivityStreamsAttributedToPropertyIterator) GetType() vocab.Type {
	return i.it.GetType()
}

func (i iter_ActivityStreamsAttributedToPropertyIterator) Next() Iter {
	return iter_ActivityStreamsAttributedToPropertyIterator{
		p:  i.p,
		it: i.it.Next(),
	}
}

func (i iter_ActivityStreamsAttributedToPropertyIterator) End() Iter {
	return iter_ActivityStreamsAttributedToPropertyIterator{
		p:  i.p,
		it: i.p.End(),
	}
}

type iter_ActivityStreamsAudiencePropertyIterator struct {
	p  vocab.ActivityStreamsAudienceProperty
	it vocab.ActivityStreamsAudiencePropertyIterator
}

func (i iter_ActivityStreamsAudiencePropertyIterator) HasAny() bool {
	return i.it.HasAny()
}

func (i iter_ActivityStreamsAudiencePropertyIterator) IsIRI() bool {
	return i.it.IsIRI()
}

func (i iter_ActivityStreamsAudiencePropertyIterator) GetIRI() *url.URL {
	return i.it.GetIRI()
}

func (i iter_ActivityStreamsAudiencePropertyIterator) GetType() vocab.Type {
	return i.it.GetType()
}

func (i iter_ActivityStreamsAudiencePropertyIterator) Next() Iter {
	return iter_ActivityStreamsAudiencePropertyIterator{
		p:  i.p,
		it: i.it.Next(),
	}
}

func (i iter_ActivityStreamsAudiencePropertyIterator) End() Iter {
	return iter_ActivityStreamsAudiencePropertyIterator{
		p:  i.p,
		it: i.p.End(),
	}
}

type iter_ActivityStreamsBccPropertyIterator struct {
	p  vocab.ActivityStreamsBccProperty
	it vocab.ActivityStreamsBccPropertyIterator
}

func (i iter_ActivityStreamsBccPropertyIterator) HasAny() bool {
	return i.it.HasAny()
}

func (i iter_ActivityStreamsBccPropertyIterator) IsIRI() bool {
	return i.it.IsIRI()
}

func (i iter_ActivityStreamsBccPropertyIterator) GetIRI() *url.URL {
	return i.it.GetIRI()
}

func (i iter_ActivityStreamsBccPropertyIterator) GetType() vocab.Type {
	return i.it.GetType()
}

func (i iter_ActivityStreamsBccPropertyIterator) Next() Iter {
	return iter_ActivityStreamsBccPropertyIterator{
		p:  i.p,
		it: i.it.Next(),
	}
}

func (i iter_ActivityStreamsBccPropertyIterator) End() Iter {
	return iter_ActivityStreamsBccPropertyIterator{
		p:  i.p,
		it: i.p.End(),
	}
}

type iter_ActivityStreamsBtoPropertyIterator struct {
	p  vocab.ActivityStreamsBtoProperty
	it vocab.ActivityStreamsBtoPropertyIterator
}

func (i iter_ActivityStreamsBtoPropertyIterator) HasAny() bool {
	return i.it.HasAny()
}

func (i iter_ActivityStreamsBtoPropertyIterator) IsIRI() bool {
	return i.it.IsIRI()
}

func (i iter_ActivityStreamsBtoPropertyIterator) GetIRI() *url.URL {
	return i.it.GetIRI()
}

func (i iter_ActivityStreamsBtoPropertyIterator) GetType() vocab.Type {
	return i.it.GetType()
}

func (i iter_ActivityStreamsBtoPropertyIterator) Next() Iter {
	return iter_ActivityStreamsBtoPropertyIterator{
		p:  i.p,
		it: i.it.Next(),
	}
}

func (i iter_ActivityStreamsBtoPropertyIterator) End() Iter {
	return iter_ActivityStreamsBtoPropertyIterator{
		p:  i.p,
		it: i.p.End(),
	}
}

type iter_ActivityStreamsCcPropertyIterator struct {
	p  vocab.ActivityStreamsCcProperty
	it vocab.ActivityStreamsCcPropertyIterator
}

func (i iter_ActivityStreamsCcPropertyIterator) HasAny() bool {
	return i.it.HasAny()
}

func (i iter_ActivityStreamsCcPropertyIterator) IsIRI() bool {
	return i.it.IsIRI()
}

func (i iter_ActivityStreamsCcPropertyIterator) GetIRI() *url.URL {
	return i.it.GetIRI()
}

func (i iter_ActivityStreamsCcPropertyIterator) GetType() vocab.Type {
	return i.it.GetType()
}

func (i iter_ActivityStreamsCcPropertyIterator) Next() Iter {
	return iter_ActivityStreamsCcPropertyIterator{
		p:  i.p,
		it: i.it.Next(),
	}
}

func (i iter_ActivityStreamsCcPropertyIterator) End() Iter {
	return iter_ActivityStreamsCcPropertyIterator{
		p:  i.p,
		it: i.p.End(),
	}
}

type iter_ActivityStreamsClosedPropertyIterator struct {
	p  vocab.ActivityStreamsClosedProperty
	it vocab.ActivityStreamsClosedPropertyIterator
}

func (i iter_ActivityStreamsClosedPropertyIterator) HasAny() bool {
	return i.it.HasAny()
}

func (i iter_ActivityStreamsClosedPropertyIterator) IsIRI() bool {
	return i.it.IsIRI()
}

func (i iter_ActivityStreamsClosedPropertyIterator) GetIRI() *url.URL {
	return i.it.GetIRI()
}

func (i iter_ActivityStreamsClosedPropertyIterator) GetType() vocab.Type {
	return i.it.GetType()
}

func (i iter_ActivityStreamsClosedPropertyIterator) Next() Iter {
	return iter_ActivityStreamsClosedPropertyIterator{
		p:  i.p,
		it: i.it.Next(),
	}
}

func (i iter_ActivityStreamsClosedPropertyIterator) End() Iter {
	return iter_ActivityStreamsClosedPropertyIterator{
		p:  i.p,
		it: i.p.End(),
	}
}

type iter_ActivityStreamsContextPropertyIterator struct {
	p  vocab.ActivityStreamsContextProperty
	it vocab.ActivityStreamsContextPropertyIterator
}

func (i iter_ActivityStreamsContextPropertyIterator) HasAny() bool {
	return i.it.HasAny()
}

func (i iter_ActivityStreamsContextPropertyIterator) IsIRI() bool {
	return i.it.IsIRI()
}

func (i iter_ActivityStreamsContextPropertyIterator) GetIRI() *url.URL {
	return i.it.GetIRI()
}

func (i iter_ActivityStreamsContextPropertyIterator) GetType() vocab.Type {
	return i.it.GetType()
}

func (i iter_ActivityStreamsContextPropertyIterator) Next() Iter {
	return iter_ActivityStreamsContextPropertyIterator{
		p:  i.p,
		it: i.it.Next(),
	}
}

func (i iter_ActivityStreamsContextPropertyIterator) End() Iter {
	return iter_ActivityStreamsContextPropertyIterator{
		p:  i.p,
		it: i.p.End(),
	}
}

type iter_ActivityStreamsFormerTypePropertyIterator struct {
	p  vocab.ActivityStreamsFormerTypeProperty
	it vocab.ActivityStreamsFormerTypePropertyIterator
}

func (i iter_ActivityStreamsFormerTypePropertyIterator) HasAny() bool {
	return i.it.HasAny()
}

func (i iter_ActivityStreamsFormerTypePropertyIterator) IsIRI() bool {
	return i.it.IsIRI()
}

func (i iter_ActivityStreamsFormerTypePropertyIterator) GetIRI() *url.URL {
	return i.it.GetIRI()
}

func (i iter_ActivityStreamsFormerTypePropertyIterator) GetType() vocab.Type {
	return i.it.GetType()
}

func (i iter_ActivityStreamsFormerTypePropertyIterator) Next() Iter {
	return iter_ActivityStreamsFormerTypePropertyIterator{
		p:  i.p,
		it: i.it.Next(),
	}
}

func (i iter_ActivityStreamsFormerTypePropertyIterator) End() Iter {
	return iter_ActivityStreamsFormerTypePropertyIterator{
		p:  i.p,
		it: i.p.End(),
	}
}

type iter_ActivityStreamsGeneratorPropertyIterator struct {
	p  vocab.ActivityStreamsGeneratorProperty
	it vocab.ActivityStreamsGeneratorPropertyIterator
}

func (i iter_ActivityStreamsGeneratorPropertyIterator) HasAny() bool {
	return i.it.HasAny()
}

func (i iter_ActivityStreamsGeneratorPropertyIterator) IsIRI() bool {
	return i.it.IsIRI()
}

func (i iter_ActivityStreamsGeneratorPropertyIterator) GetIRI() *url.URL {
	return i.it.GetIRI()
}

func (i iter_ActivityStreamsGeneratorPropertyIterator) GetType() vocab.Type {
	return i.it.GetType()
}

func (i iter_ActivityStreamsGeneratorPropertyIterator) Next() Iter {
	return iter_ActivityStreamsGeneratorPropertyIterator{
		p:  i.p,
		it: i.it.Next(),
	}
}

func (i iter_ActivityStreamsGeneratorPropertyIterator) End() Iter {
	return iter_ActivityStreamsGeneratorPropertyIterator{
		p:  i.p,
		it: i.p.End(),
	}
}

type iter_ActivityStreamsIconPropertyIterator struct {
	p  vocab.ActivityStreamsIconProperty
	it vocab.ActivityStreamsIconPropertyIterator
}

func (i iter_ActivityStreamsIconPropertyIterator) HasAny() bool {
	return i.it.HasAny()
}

func (i iter_ActivityStreamsIconPropertyIterator) IsIRI() bool {
	return i.it.IsIRI()
}

func (i iter_ActivityStreamsIconPropertyIterator) GetIRI() *url.URL {
	return i.it.GetIRI()
}

func (i iter_ActivityStreamsIconPropertyIterator) GetType() vocab.Type {
	return i.it.GetType()
}

func (i iter_ActivityStreamsIconPropertyIterator) Next() Iter {
	return iter_ActivityStreamsIconPropertyIterator{
		p:  i.p,
		it: i.it.Next(),
	}
}

func (i iter_ActivityStreamsIconPropertyIterator) End() Iter {
	return iter_ActivityStreamsIconPropertyIterator{
		p:  i.p,
		it: i.p.End(),
	}
}

type iter_ActivityStreamsImagePropertyIterator struct {
	p  vocab.ActivityStreamsImageProperty
	it vocab.ActivityStreamsImagePropertyIterator
}

func (i iter_ActivityStreamsImagePropertyIterator) HasAny() bool {
	return i.it.HasAny()
}

func (i iter_ActivityStreamsImagePropertyIterator) IsIRI() bool {
	return i.it.IsIRI()
}

func (i iter_ActivityStreamsImagePropertyIterator) GetIRI() *url.URL {
	return i.it.GetIRI()
}

func (i iter_ActivityStreamsImagePropertyIterator) GetType() vocab.Type {
	return i.it.GetType()
}

func (i iter_ActivityStreamsImagePropertyIterator) Next() Iter {
	return iter_ActivityStreamsImagePropertyIterator{
		p:  i.p,
		it: i.it.Next(),
	}
}

func (i iter_ActivityStreamsImagePropertyIterator) End() Iter {
	return iter_ActivityStreamsImagePropertyIterator{
		p:  i.p,
		it: i.p.End(),
	}
}

type iter_ActivityStreamsInReplyToPropertyIterator struct {
	p  vocab.ActivityStreamsInReplyToProperty
	it vocab.ActivityStreamsInReplyToPropertyIterator
}

func (i iter_ActivityStreamsInReplyToPropertyIterator) HasAny() bool {
	return i.it.HasAny()
}

func (i iter_ActivityStreamsInReplyToPropertyIterator) IsIRI() bool {
	return i.it.IsIRI()
}

func (i iter_ActivityStreamsInReplyToPropertyIterator) GetIRI() *url.URL {
	return i.it.GetIRI()
}

func (i iter_ActivityStreamsInReplyToPropertyIterator) GetType() vocab.Type {
	return i.it.GetType()
}

func (i iter_ActivityStreamsInReplyToPropertyIterator) Next() Iter {
	return iter_ActivityStreamsInReplyToPropertyIterator{
		p:  i.p,
		it: i.it.Next(),
	}
}

func (i iter_ActivityStreamsInReplyToPropertyIterator) End() Iter {
	return iter_ActivityStreamsInReplyToPropertyIterator{
		p:  i.p,
		it: i.p.End(),
	}
}

type iter_ActivityStreamsInstrumentPropertyIterator struct {
	p  vocab.ActivityStreamsInstrumentProperty
	it vocab.ActivityStreamsInstrumentPropertyIterator
}

func (i iter_ActivityStreamsInstrumentPropertyIterator) HasAny() bool {
	return i.it.HasAny()
}

func (i iter_ActivityStreamsInstrumentPropertyIterator) IsIRI() bool {
	return i.it.IsIRI()
}

func (i iter_ActivityStreamsInstrumentPropertyIterator) GetIRI() *url.URL {
	return i.it.GetIRI()
}

func (i iter_ActivityStreamsInstrumentPropertyIterator) GetType() vocab.Type {
	return i.it.GetType()
}

func (i iter_ActivityStreamsInstrumentPropertyIterator) Next() Iter {
	return iter_ActivityStreamsInstrumentPropertyIterator{
		p:  i.p,
		it: i.it.Next(),
	}
}

func (i iter_ActivityStreamsInstrumentPropertyIterator) End() Iter {
	return iter_ActivityStreamsInstrumentPropertyIterator{
		p:  i.p,
		it: i.p.End(),
	}
}

type iter_ActivityStreamsItemsPropertyIterator struct {
	p  vocab.ActivityStreamsItemsProperty
	it vocab.ActivityStreamsItemsPropertyIterator
}

func (i iter_ActivityStreamsItemsPropertyIterator) HasAny() bool {
	return i.it.HasAny()
}

func (i iter_ActivityStreamsItemsPropertyIterator) IsIRI() bool {
	return i.it.IsIRI()
}

func (i iter_ActivityStreamsItemsPropertyIterator) GetIRI() *url.URL {
	return i.it.GetIRI()
}

func (i iter_ActivityStreamsItemsPropertyIterator) GetType() vocab.Type {
	return i.it.GetType()
}

func (i iter_ActivityStreamsItemsPropertyIterator) Next() Iter {
	return iter_ActivityStreamsItemsPropertyIterator{
		p:  i.p,
		it: i.it.Next(),
	}
}

func (i iter_ActivityStreamsItemsPropertyIterator) End() Iter {
	return iter_ActivityStreamsItemsPropertyIterator{
		p:  i.p,
		it: i.p.End(),
	}
}

type iter_ActivityStreamsLocationPropertyIterator struct {
	p  vocab.ActivityStreamsLocationProperty
	it vocab.ActivityStreamsLocationPropertyIterator
}

func (i iter_ActivityStreamsLocationPropertyIterator) HasAny() bool {
	return i.it.HasAny()
}

func (i iter_ActivityStreamsLocationPropertyIterator) IsIRI() bool {
	return i.it.IsIRI()
}

func (i iter_ActivityStreamsLocationPropertyIterator) GetIRI() *url.URL {
	return i.it.GetIRI()
}

func (i iter_ActivityStreamsLocationPropertyIterator) GetType() vocab.Type {
	return i.it.GetType()
}

func (i iter_ActivityStreamsLocationPropertyIterator) Next() Iter {
	return iter_ActivityStreamsLocationPropertyIterator{
		p:  i.p,
		it: i.it.Next(),
	}
}

func (i iter_ActivityStreamsLocationPropertyIterator) End() Iter {
	return iter_ActivityStreamsLocationPropertyIterator{
		p:  i.p,
		it: i.p.End(),
	}
}

type iter_ActivityStreamsObjectPropertyIterator struct {
	p  vocab.ActivityStreamsObjectProperty
	it vocab.ActivityStreamsObjectPropertyIterator
}

func (i iter_ActivityStreamsObjectPropertyIterator) HasAny() bool {
	return i.it.HasAny()
}

func (i iter_ActivityStreamsObjectPropertyIterator) IsIRI() bool {
	return i.it.IsIRI()
}

func (i iter_ActivityStreamsObjectPropertyIterator) GetIRI() *url.URL {
	return i.it.GetIRI()
}

func (i iter_ActivityStreamsObjectPropertyIterator) GetType() vocab.Type {
	return i.it.GetType()
}

func (i iter_ActivityStreamsObjectPropertyIterator) Next() Iter {
	return iter_ActivityStreamsObjectPropertyIterator{
		p:  i.p,
		it: i.it.Next(),
	}
}

func (i iter_ActivityStreamsObjectPropertyIterator) End() Iter {
	return iter_ActivityStreamsObjectPropertyIterator{
		p:  i.p,
		it: i.p.End(),
	}
}

type iter_ActivityStreamsOneOfPropertyIterator struct {
	p  vocab.ActivityStreamsOneOfProperty
	it vocab.ActivityStreamsOneOfPropertyIterator
}

func (i iter_ActivityStreamsOneOfPropertyIterator) HasAny() bool {
	return i.it.HasAny()
}

func (i iter_ActivityStreamsOneOfPropertyIterator) IsIRI() bool {
	return i.it.IsIRI()
}

func (i iter_ActivityStreamsOneOfPropertyIterator) GetIRI() *url.URL {
	return i.it.GetIRI()
}

func (i iter_ActivityStreamsOneOfPropertyIterator) GetType() vocab.Type {
	return i.it.GetType()
}

func (i iter_ActivityStreamsOneOfPropertyIterator) Next() Iter {
	return iter_ActivityStreamsOneOfPropertyIterator{
		p:  i.p,
		it: i.it.Next(),
	}
}

func (i iter_ActivityStreamsOneOfPropertyIterator) End() Iter {
	return iter_ActivityStreamsOneOfPropertyIterator{
		p:  i.p,
		it: i.p.End(),
	}
}

type iter_ActivityStreamsOrderedItemsPropertyIterator struct {
	p  vocab.ActivityStreamsOrderedItemsProperty
	it vocab.ActivityStreamsOrderedItemsPropertyIterator
}

func (i iter_ActivityStreamsOrderedItemsPropertyIterator) HasAny() bool {
	return i.it.HasAny()
}

func (i iter_ActivityStreamsOrderedItemsPropertyIterator) IsIRI() bool {
	return i.it.IsIRI()
}

func (i iter_ActivityStreamsOrderedItemsPropertyIterator) GetIRI() *url.URL {
	return i.it.GetIRI()
}

func (i iter_ActivityStreamsOrderedItemsPropertyIterator) GetType() vocab.Type {
	return i.it.GetType()
}

func (i iter_ActivityStreamsOrderedItemsPropertyIterator) Next() Iter {
	return iter_ActivityStreamsOrderedItemsPropertyIterator{
		p:  i.p,
		it: i.it.Next(),
	}
}

func (i iter_ActivityStreamsOrderedItemsPropertyIterator) End() Iter {
	return iter_ActivityStreamsOrderedItemsPropertyIterator{
		p:  i.p,
		it: i.p.End(),
	}
}

type iter_ActivityStreamsOriginPropertyIterator struct {
	p  vocab.ActivityStreamsOriginProperty
	it vocab.ActivityStreamsOriginPropertyIterator
}

func (i iter_ActivityStreamsOriginPropertyIterator) HasAny() bool {
	return i.it.HasAny()
}

func (i iter_ActivityStreamsOriginPropertyIterator) IsIRI() bool {
	return i.it.IsIRI()
}

func (i iter_ActivityStreamsOriginPropertyIterator) GetIRI() *url.URL {
	return i.it.GetIRI()
}

func (i iter_ActivityStreamsOriginPropertyIterator) GetType() vocab.Type {
	return i.it.GetType()
}

func (i iter_ActivityStreamsOriginPropertyIterator) Next() Iter {
	return iter_ActivityStreamsOriginPropertyIterator{
		p:  i.p,
		it: i.it.Next(),
	}
}

func (i iter_ActivityStreamsOriginPropertyIterator) End() Iter {
	return iter_ActivityStreamsOriginPropertyIterator{
		p:  i.p,
		it: i.p.End(),
	}
}

type iter_ActivityStreamsPreviewPropertyIterator struct {
	p  vocab.ActivityStreamsPreviewProperty
	it vocab.ActivityStreamsPreviewPropertyIterator
}

func (i iter_ActivityStreamsPreviewPropertyIterator) HasAny() bool {
	return i.it.HasAny()
}

func (i iter_ActivityStreamsPreviewPropertyIterator) IsIRI() bool {
	return i.it.IsIRI()
}

func (i iter_ActivityStreamsPreviewPropertyIterator) GetIRI() *url.URL {
	return i.it.GetIRI()
}

func (i iter_ActivityStreamsPreviewPropertyIterator) GetType() vocab.Type {
	return i.it.GetType()
}

func (i iter_ActivityStreamsPreviewPropertyIterator) Next() Iter {
	return iter_ActivityStreamsPreviewPropertyIterator{
		p:  i.p,
		it: i.it.Next(),
	}
}

func (i iter_ActivityStreamsPreviewPropertyIterator) End() Iter {
	return iter_ActivityStreamsPreviewPropertyIterator{
		p:  i.p,
		it: i.p.End(),
	}
}

type iter_ActivityStreamsRelationshipPropertyIterator struct {
	p  vocab.ActivityStreamsRelationshipProperty
	it vocab.ActivityStreamsRelationshipPropertyIterator
}

func (i iter_ActivityStreamsRelationshipPropertyIterator) HasAny() bool {
	return i.it.HasAny()
}

func (i iter_ActivityStreamsRelationshipPropertyIterator) IsIRI() bool {
	return i.it.IsIRI()
}

func (i iter_ActivityStreamsRelationshipPropertyIterator) GetIRI() *url.URL {
	return i.it.GetIRI()
}

func (i iter_ActivityStreamsRelationshipPropertyIterator) GetType() vocab.Type {
	return i.it.GetType()
}

func (i iter_ActivityStreamsRelationshipPropertyIterator) Next() Iter {
	return iter_ActivityStreamsRelationshipPropertyIterator{
		p:  i.p,
		it: i.it.Next(),
	}
}

func (i iter_ActivityStreamsRelationshipPropertyIterator) End() Iter {
	return iter_ActivityStreamsRelationshipPropertyIterator{
		p:  i.p,
		it: i.p.End(),
	}
}

type iter_ActivityStreamsResultPropertyIterator struct {
	p  vocab.ActivityStreamsResultProperty
	it vocab.ActivityStreamsResultPropertyIterator
}

func (i iter_ActivityStreamsResultPropertyIterator) HasAny() bool {
	return i.it.HasAny()
}

func (i iter_ActivityStreamsResultPropertyIterator) IsIRI() bool {
	return i.it.IsIRI()
}

func (i iter_ActivityStreamsResultPropertyIterator) GetIRI() *url.URL {
	return i.it.GetIRI()
}

func (i iter_ActivityStreamsResultPropertyIterator) GetType() vocab.Type {
	return i.it.GetType()
}

func (i iter_ActivityStreamsResultPropertyIterator) Next() Iter {
	return iter_ActivityStreamsResultPropertyIterator{
		p:  i.p,
		it: i.it.Next(),
	}
}

func (i iter_ActivityStreamsResultPropertyIterator) End() Iter {
	return iter_ActivityStreamsResultPropertyIterator{
		p:  i.p,
		it: i.p.End(),
	}
}

type iter_ActivityStreamsStreamsPropertyIterator struct {
	p  vocab.ActivityStreamsStreamsProperty
	it vocab.ActivityStreamsStreamsPropertyIterator
}

func (i iter_ActivityStreamsStreamsPropertyIterator) HasAny() bool {
	return i.it.HasAny()
}

func (i iter_ActivityStreamsStreamsPropertyIterator) IsIRI() bool {
	return i.it.IsIRI()
}

func (i iter_ActivityStreamsStreamsPropertyIterator) GetIRI() *url.URL {
	return i.it.GetIRI()
}

func (i iter_ActivityStreamsStreamsPropertyIterator) GetType() vocab.Type {
	return i.it.GetType()
}

func (i iter_ActivityStreamsStreamsPropertyIterator) Next() Iter {
	return iter_ActivityStreamsStreamsPropertyIterator{
		p:  i.p,
		it: i.it.Next(),
	}
}

func (i iter_ActivityStreamsStreamsPropertyIterator) End() Iter {
	return iter_ActivityStreamsStreamsPropertyIterator{
		p:  i.p,
		it: i.p.End(),
	}
}

type iter_ActivityStreamsTagPropertyIterator struct {
	p  vocab.ActivityStreamsTagProperty
	it vocab.ActivityStreamsTagPropertyIterator
}

func (i iter_ActivityStreamsTagPropertyIterator) HasAny() bool {
	return i.it.HasAny()
}

func (i iter_ActivityStreamsTagPropertyIterator) IsIRI() bool {
	return i.it.IsIRI()
}

func (i iter_ActivityStreamsTagPropertyIterator) GetIRI() *url.URL {
	return i.it.GetIRI()
}

func (i iter_ActivityStreamsTagPropertyIterator) GetType() vocab.Type {
	return i.it.GetType()
}

func (i iter_ActivityStreamsTagPropertyIterator) Next() Iter {
	return iter_ActivityStreamsTagPropertyIterator{
		p:  i.p,
		it: i.it.Next(),
	}
}

func (i iter_ActivityStreamsTagPropertyIterator) End() Iter {
	return iter_ActivityStreamsTagPropertyIterator{
		p:  i.p,
		it: i.p.End(),
	}
}

type iter_ActivityStreamsTargetPropertyIterator struct {
	p  vocab.ActivityStreamsTargetProperty
	it vocab.ActivityStreamsTargetPropertyIterator
}

func (i iter_ActivityStreamsTargetPropertyIterator) HasAny() bool {
	return i.it.HasAny()
}

func (i iter_ActivityStreamsTargetPropertyIterator) IsIRI() bool {
	return i.it.IsIRI()
}

func (i iter_ActivityStreamsTargetPropertyIterator) GetIRI() *url.URL {
	return i.it.GetIRI()
}

func (i iter_ActivityStreamsTargetPropertyIterator) GetType() vocab.Type {
	return i.it.GetType()
}

func (i iter_ActivityStreamsTargetPropertyIterator) Next() Iter {
	return iter_ActivityStreamsTargetPropertyIterator{
		p:  i.p,
		it: i.it.Next(),
	}
}

func (i iter_ActivityStreamsTargetPropertyIterator) End() Iter {
	return iter_ActivityStreamsTargetPropertyIterator{
		p:  i.p,
		it: i.p.End(),
	}
}

type iter_ActivityStreamsToPropertyIterator struct {
	p  vocab.ActivityStreamsToProperty
	it vocab.ActivityStreamsToPropertyIterator
}

func (i iter_ActivityStreamsToPropertyIterator) HasAny() bool {
	return i.it.HasAny()
}

func (i iter_ActivityStreamsToPropertyIterator) IsIRI() bool {
	return i.it.IsIRI()
}

func (i iter_ActivityStreamsToPropertyIterator) GetIRI() *url.URL {
	return i.it.GetIRI()
}

func (i iter_ActivityStreamsToPropertyIterator) GetType() vocab.Type {
	return i.it.GetType()
}

func (i iter_ActivityStreamsToPropertyIterator) Next() Iter {
	return iter_ActivityStreamsToPropertyIterator{
		p:  i.p,
		it: i.it.Next(),
	}
}

func (i iter_ActivityStreamsToPropertyIterator) End() Iter {
	return iter_ActivityStreamsToPropertyIterator{
		p:  i.p,
		it: i.p.End(),
	}
}

type iter_ActivityStreamsUrlPropertyIterator struct {
	p  vocab.ActivityStreamsUrlProperty
	it vocab.ActivityStreamsUrlPropertyIterator
}

func (i iter_ActivityStreamsUrlPropertyIterator) HasAny() bool {
	return i.it.HasAny()
}

func (i iter_ActivityStreamsUrlPropertyIterator) IsIRI() bool {
	return i.it.IsIRI()
}

func (i iter_ActivityStreamsUrlPropertyIterator) GetIRI() *url.URL {
	return i.it.GetIRI()
}

func (i iter_ActivityStreamsUrlPropertyIterator) GetType() vocab.Type {
	return i.it.GetType()
}

func (i iter_ActivityStreamsUrlPropertyIterator) Next() Iter {
	return iter_ActivityStreamsUrlPropertyIterator{
		p:  i.p,
		it: i.it.Next(),
	}
}

func (i iter_ActivityStreamsUrlPropertyIterator) End() Iter {
	return iter_ActivityStreamsUrlPropertyIterator{
		p:  i.p,
		it: i.p.End(),
	}
}

type iter_W3IDSecurityV1PublicKeyPropertyIterator struct {
	p  vocab.W3IDSecurityV1PublicKeyProperty
	it vocab.W3IDSecurityV1PublicKeyPropertyIterator
}

func (i iter_W3IDSecurityV1PublicKeyPropertyIterator) HasAny() bool {
	return i.it.HasAny()
}

func (i iter_W3IDSecurityV1PublicKeyPropertyIterator) IsIRI() bool {
	return i.it.IsIRI()
}

func (i iter_W3IDSecurityV1PublicKeyPropertyIterator) GetIRI() *url.URL {
	return i.it.GetIRI()
}

func (i iter_W3IDSecurityV1PublicKeyPropertyIterator) GetType() vocab.Type {
	return i.it.GetType()
}

func (i iter_W3IDSecurityV1PublicKeyPropertyIterator) Next() Iter {
	return iter_W3IDSecurityV1PublicKeyPropertyIterator{
		p:  i.p,
		it: i.it.Next(),
	}
}

func (i iter_W3IDSecurityV1PublicKeyPropertyIterator) End() Iter {
	return iter_W3IDSecurityV1PublicKeyPropertyIterator{
		p:  i.p,
		it: i.p.End(),
	}
}

#! /usr/bin/env python3


'''
iterator_impl.gen.py

This is a script to generate the iterator_impl.go file which is used
to make the terrible iterator API usable. Why am I doing this to
myself?
'''


iterator_types = (
    'ActivityStreamsActorPropertyIterator',
    'ActivityStreamsAnyOfPropertyIterator',
    'ActivityStreamsAttachmentPropertyIterator',
    'ActivityStreamsAttributedToPropertyIterator',
    'ActivityStreamsAudiencePropertyIterator',
    'ActivityStreamsBccPropertyIterator',
    'ActivityStreamsBtoPropertyIterator',
    'ActivityStreamsCcPropertyIterator',
    'ActivityStreamsClosedPropertyIterator',
    'ActivityStreamsContextPropertyIterator',
    'ActivityStreamsFormerTypePropertyIterator',
    'ActivityStreamsGeneratorPropertyIterator',
    'ActivityStreamsIconPropertyIterator',
    'ActivityStreamsImagePropertyIterator',
    'ActivityStreamsInReplyToPropertyIterator',
    'ActivityStreamsInstrumentPropertyIterator',
    'ActivityStreamsItemsPropertyIterator',
    'ActivityStreamsLocationPropertyIterator',
    'ActivityStreamsObjectPropertyIterator',
    'ActivityStreamsOneOfPropertyIterator',
    'ActivityStreamsOrderedItemsPropertyIterator',
    'ActivityStreamsOriginPropertyIterator',
    'ActivityStreamsPreviewPropertyIterator',
    'ActivityStreamsRelationshipPropertyIterator',
    'ActivityStreamsResultPropertyIterator',
    'ActivityStreamsStreamsPropertyIterator',
    'ActivityStreamsTagPropertyIterator',
    'ActivityStreamsTargetPropertyIterator',
    'ActivityStreamsToPropertyIterator',
    'ActivityStreamsUrlPropertyIterator',
    'W3IDSecurityV1PublicKeyPropertyIterator',
)


not_supported = (
    'ActivityStreamsContentPropertyIterator',
    'ActivityStreamsNamePropertyIterator',
    'ActivityStreamsRelPropertyIterator',
    'ActivityStreamsSummaryPropertyIterator',
    'JSONLDTypePropertyIterator',
)


wrapper_template = '''
type iter_$ITER_TYPE struct {
        p vocab.$BASE_TYPE
        it vocab.$ITER_TYPE
}

func (i iter_$ITER_TYPE) HasAny() bool {
        return i.it.HasAny()
}

func (i iter_$ITER_TYPE) IsIRI() bool {
        return i.it.IsIRI()
}

func (i iter_$ITER_TYPE) GetIRI() *url.URL {
        return i.it.GetIRI()
}

func (i iter_$ITER_TYPE) GetType() vocab.Type {
        return i.it.GetType()
}

func (i iter_$ITER_TYPE) Next() Iter {
        return iter_$ITER_TYPE{
                p: i.p,
                it: i.it.Next(),
        }
}

func (i iter_$ITER_TYPE) End() Iter {
        return iter_$ITER_TYPE{
                p: i.p,
                it: i.p.End(),
        }
}'''


def emit_wrapper(itername: str):
    basename = itername.replace('Iterator', '')

    wrapper = wrapper_template
    wrapper = wrapper.replace('$BASE_TYPE', basename)
    wrapper = wrapper.replace('$ITER_TYPE', itername)

    print(wrapper)


def emit_wrappers():
    for itername in iterator_types:
        emit_wrapper(itername)


def emit_constructor_for(itername: str):
    basename = itername.replace('Iterator', '')

    print('	case vocab.%s:' % basename)
    print('		return iter_%s{' % itername)
    print('			p: v,')
    print('			it: v.Begin(),')
    print('		}, nil')
    print()


def emit_constructor():
    print('func begin(iterable interface{}) (Iter, error) {')
    print('if iterable == nil {')
    print('	return nil, errors.New("nil argument")')
    print('}')
    print()
    print('	switch v := iterable.(type) {')

    for itername in iterator_types:
        emit_constructor_for(itername)

    print('	default:')
    print('		return nil, fmt.Errorf("type=%T not supported", iterable)')
    print('	}')
    print('}')


def emit_header():
    print('package fetch')
    print()
    print('// AUTO GENERATED')
    print('// see iterator_impl.gen.py for details')
    print()
    print('import (')
    print('	"errors"')
    print('	"fmt"')
    print('	"github.com/go-fed/activity/streams/vocab"')
    print('	"net/url"')
    print(')')
    print()


def main():
    emit_header()
    emit_constructor()
    emit_wrappers()


if __name__ == '__main__':
    try:
        main()
    except (BrokenPipeError, KeyboardInterrupt, SystemExit):
        pass

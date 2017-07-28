import _ from 'lodash'
import $ from 'jquery'
import Vue from 'vue'
import BaseWidget from 'components/widget/widget.vue'

const WIDGET_SELECTOR = '[data-sp_widget]'

function getElements () {
  return $(WIDGET_SELECTOR).get()
}

function initWidgetElement (element) {
  const ds = element.dataset || {}
  if (!ds.sp_widget) throw new Error('Widget needs an endpoint.')

  let props = JSON.parse(ds.parameters)
  props = _.mapKeys(props, (value, key) => _.camelCase(key))
  props.endpoint = ds.sp_widget

  return new Vue({
    el: element,
    render (createElement) {
      return createElement(BaseWidget, {
        props
      })
    }
  })
}

export default function initWidgets () {
  getElements().forEach(initWidgetElement)
}

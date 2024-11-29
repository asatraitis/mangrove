import { LitElement, css, html } from 'lit'
import { customElement, property } from 'lit/decorators.js'

import '@material/web/textfield/outlined-text-field';

/**
 * An example element.
 *
 * @slot - This element has a slot
 * @csspart button - The button
 */
@customElement('mg-app')
export class App extends LitElement {
  /**
   * Copy for the read the docs hint.
   */
  @property()
  docsHint = 'Click on the Vite and Lit logos to learn more'

  /**
   * The number of times the button has been clicked.
   */
  @property({ type: Number })
  count = 0

  render() {
    return html`
      <div>
        <md-outlined-text-field label="Outlined" value="Value"></md-outlined-text-field>
        ${this.count}
      </div>
    `
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'mg-app': App
  }
}

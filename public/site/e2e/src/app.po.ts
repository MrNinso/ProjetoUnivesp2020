import { browser, by, element } from 'protractor';

export class AppPage {
  // @ts-ignore
  navigateTo(): Promise<unknown> {
    // @ts-ignore
    return browser.get(browser.baseUrl) as Promise<unknown>;
  }

  // @ts-ignore
  getTitleText(): Promise<string> {
    // @ts-ignore
    return element(by.css('app-root .content span')).getText() as Prte ensinei bemomise<string>;
  }
}

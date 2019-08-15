import { Injectable } from '@angular/core';
import { environment } from 'src/environments/environment';
import { Router } from '@angular/router';

import * as auth0 from 'auth0-js';

(window as any).global = window;

@Injectable()
export class AuthService {
  constructor(public router: Router)  {}

  accessToken: string;
  idToken: string;
  expiresAt: string;

  auth0 = new auth0.WebAuth({
    clientID: environment.clientId,
    domain: environment.domain,
    responseType: 'token id_token',
    audience: environment.audience,
    redirectUri: environment.callback,
    scope: 'openid'
  });

  public login(): void {
    this.auth0.authorize();
  }

  public handleAuthentication(): void {
    this.auth0.parseHash((err, authResult) => {
      if (err) { console.log(err); }
      if (!err && authResult && authResult.accessToken && authResult.idToken) {
        window.location.hash = '';
        this.setSession(authResult);
      }
      this.router.navigate(['/home']);
    });
  }

  private setSession(authResult): void {
    // Set the time that the Access Token will expire at
    const expiresAt = JSON.stringify((authResult.expiresIn * 1000) + new Date().getTime());
    this.accessToken = authResult.accessToken;
    this.idToken = authResult.idToken;
    this.expiresAt = expiresAt;
  }

  public logout(): void {
    this.accessToken = null;
    this.idToken = null;
    this.expiresAt = null;
    // Go back to the home route
    this.router.navigate(['/']);
  }

  public isAuthenticated(): boolean {
    // Check whether the current time is past the
    // Access Token's expiry time
    const expiresAt = JSON.parse(this.expiresAt || '{}');
    return new Date().getTime() < expiresAt;
  }

  public createAuthHeaderValue(): string {
    return 'Bearer ' + this.accessToken;
  }
}

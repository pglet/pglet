import React from 'react'
import { Dialog, DialogType, IDialogProps, Image, Stack, Text, DefaultButton, Checkbox } from '@fluentui/react';
import { ILoginProps } from './Control.types'
import pglet_logo from '../assets/img/pglet-logo-no-text.svg'
import microsoft_logo from '../assets/img/microsoft-logo.svg'
import google_logo from '../assets/img/google-logo.svg'
import github_logo from '../assets/img/github-logo.svg'

export const Login = React.memo<ILoginProps>(({loginOptions, onDismiss}) => {

    var pageUrl = encodeURIComponent(window.location.pathname);

    const [persistLogin, setPersistLogin] = React.useState<boolean>(true);

    // dialog props
    const props: IDialogProps = {
        hidden: false,
        onDismiss: onDismiss,
        modalProps: {
            topOffsetFixed: false,
            isBlocking: onDismiss ? true : false
        },
        dialogContentProps: {
            type: onDismiss ? DialogType.close : DialogType.normal,
            title: 'Sign in to Pglet',
            styles: {
                title: {
                    textAlign: onDismiss ? 'left' : 'center'
                }
            }
        },
        styles: {
            main: {
                padding: "0 0",
                ".ms-Dialog-title": {
                    padding: '16px 24px'
                }
            },
        },
    };

    const onChange = React.useCallback((ev?: React.FormEvent<HTMLElement | HTMLInputElement>, checked?: boolean): void => {
        setPersistLogin(!!checked);
      }, []);

    const getOAuthURL = (groupsEnabled:boolean): string => {
        return `?redirect_url=${pageUrl}&persist=${persistLogin ? '1' : '0'}&groups=${groupsEnabled ? '1' : '0'}`
    }

    return <Dialog {...props}>
        <Stack horizontalAlign="center" tokens={{ childrenGap: 30}}>
            <Stack horizontalAlign="center" style={{margin: '10px 0 0 0'}}>
                <Image src={pglet_logo} width={50} height={50} />
            </Stack>
            <Text variant="medium" style={{ textAlign: "center" }}>You must sign in to access this page. Please continue with one of the options below:</Text>
            <Stack tokens={{ childrenGap: 20}} horizontalAlign="stretch">
                {
                    loginOptions.gitHubEnabled &&
                        <DefaultButton href={"/api/oauth/github" + getOAuthURL(loginOptions.gitHubGroupScope)} iconProps={{
                            imageProps: {
                                src: github_logo,
                                width: 16,
                                height: 16
                            }
                        }} style={{ padding: "0 50px" }}>Sign in with GitHub</DefaultButton>
                }
                {
                    loginOptions.googleEnabled &&
                        <DefaultButton href={"/api/oauth/google" + getOAuthURL(loginOptions.googleGroupScope)} iconProps={{
                            imageProps: {
                                src: google_logo,
                                width: 16,
                                height: 16
                            }
                        }} style={{ padding: "0 50px" }}>Sign in with Google</DefaultButton>                    
                }                
                {
                    loginOptions.azureEnabled &&
                        <DefaultButton href={"/api/oauth/azure" + getOAuthURL(loginOptions.gitHubGroupScope)} iconProps={{
                            imageProps: {
                                src: microsoft_logo,
                                width: 16,
                                height: 16
                            }
                        }}>Sign in with Microsoft account</DefaultButton>                    
                }
                <Stack horizontalAlign="center">
                    <Checkbox label="Stay signed in for a week" checked={persistLogin} onChange={onChange} />
                </Stack>
            </Stack>
        </Stack>
    </Dialog>
})
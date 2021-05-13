import React from 'react'
import { Dialog, DialogType, IDialogProps, Image, Stack, Text, DefaultButton, Checkbox } from '@fluentui/react';
import { ILoginProps } from './Control.types'
import logo_light from '../assets/img/logo_light.svg'
import microsoft_logo from '../assets/img/microsoft-logo.svg'
import google_logo from '../assets/img/google-logo.svg'
import github_logo from '../assets/img/github-logo.svg'

export const Login = React.memo<ILoginProps>(({loginOptions}) => {

    var pageUrl = encodeURIComponent(window.location.pathname);

    const [persistLogin, setPersistLogin] = React.useState<boolean>(true);

    // dialog props
    const props: IDialogProps = {
        hidden: false,
        modalProps: {
            topOffsetFixed: false,
            isBlocking: false,
        },
        dialogContentProps: {
            type: DialogType.normal,
            title: <>
                <Stack horizontalAlign="center">
                    <Image src={logo_light} />
                </Stack>
            </>,
            subText: undefined,
        },
        styles: {
            main: {
                padding: "20px 0",
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
            <Text variant="xLarge">Sign in to Pglet</Text>
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
                <Checkbox label="Stay signed in for a week" checked={persistLogin} onChange={onChange} />
            </Stack>
        </Stack>
    </Dialog>
})
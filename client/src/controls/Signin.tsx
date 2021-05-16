import React from 'react'
import { Dialog, DialogType, IDialogProps, Image, Stack, Text, DefaultButton, Checkbox, useTheme } from '@fluentui/react';
import { ISigninProps } from './Control.types'
import pglet_logo from '../assets/img/pglet-logo-no-text.svg'
import microsoft_logo from '../assets/img/microsoft-logo.svg'
import google_logo from '../assets/img/google-logo.svg'
import github_logo from '../assets/img/github-logo.svg'
import github_logo_white from '../assets/img/github-logo-white.svg'

export const Signin = React.memo<ISigninProps>(({signinOptions, onDismiss}) => {

    var pageUrl = encodeURIComponent(window.location.pathname);

    const theme = useTheme();
    const [persistSignin, setPersistSignin] = React.useState<boolean>(true);

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
        setPersistSignin(!!checked);
      }, []);

    const getOAuthURL = (groupsEnabled:boolean): string => {
        return `?redirect_url=${pageUrl}&persist=${persistSignin ? '1' : '0'}&groups=${groupsEnabled ? '1' : '0'}`
    }

    return <Dialog {...props}>
        <Stack horizontalAlign="center" tokens={{ childrenGap: 30}}>
            <Stack horizontalAlign="center" style={{margin: '10px 0 0 0'}}>
                <Image src={pglet_logo} width={50} height={50} />
            </Stack>
            <Text variant="medium" style={{ textAlign: "center" }}>You must sign in to access this page. Please continue with one of the options below:</Text>
            <Stack tokens={{ childrenGap: 20}} horizontalAlign="stretch">
                {
                    signinOptions.gitHubEnabled &&
                        <DefaultButton href={"/api/oauth/github" + getOAuthURL(signinOptions.gitHubGroupScope)} iconProps={{
                            imageProps: {
                                src: theme.isInverted ? github_logo_white : github_logo,
                                width: 16,
                                height: 16
                            }
                        }} style={{ padding: "0 50px" }}>Sign in with GitHub</DefaultButton>
                }
                {
                    signinOptions.googleEnabled &&
                        <DefaultButton href={"/api/oauth/google" + getOAuthURL(signinOptions.googleGroupScope)} iconProps={{
                            imageProps: {
                                src: google_logo,
                                width: 16,
                                height: 16
                            }
                        }} style={{ padding: "0 50px" }}>Sign in with Google</DefaultButton>                    
                }                
                {
                    signinOptions.azureEnabled &&
                        <DefaultButton href={"/api/oauth/azure" + getOAuthURL(signinOptions.gitHubGroupScope)} iconProps={{
                            imageProps: {
                                src: microsoft_logo,
                                width: 16,
                                height: 16
                            }
                        }}>Sign in with Microsoft account</DefaultButton>                    
                }
                <Stack horizontalAlign="center">
                    <Checkbox label="Stay signed in for a week" checked={persistSignin} onChange={onChange} />
                </Stack>
            </Stack>
        </Stack>
    </Dialog>
})
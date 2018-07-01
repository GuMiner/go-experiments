using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Data;
using System.Drawing;
using System.Drawing.Imaging;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Windows.Forms;

namespace ReceiptEditor
{
    public partial class ImageEditForm : Form
    {
        private Action<ImageAttributes> editCallback;

        public ImageEditForm(Action<ImageAttributes> editCallback)
        {
            InitializeComponent();
            this.editCallback = editCallback;
        }

        private void brightnessTrackbar_Scroll(object sender, EventArgs e)
        {
            UpdateImageFromEdits();
        }

        private void contrastTrackBar_Scroll(object sender, EventArgs e)
        {
            UpdateImageFromEdits();
        }

        private void UpdateImageFromEdits()
        {
            ImageAttributes imageAttributes = new ImageAttributes();

            imageAttributes.SetColorMatrix(CreateColorMatrix(brightnessTrackbar.Value, contrastTrackBar.Value),
                ColorMatrixFlag.Default, ColorAdjustType.Bitmap);

            editCallback(imageAttributes);
        }

        private ColorMatrix CreateColorMatrix(int trackbarBrighness, int trackbarContrast)
        {
            float brightness = (float)(trackbarBrighness - 50) / 50.0f;
            float contrast = 1.0f + (float)(trackbarContrast - 50) / 50.0f;
            return new ColorMatrix(
                new[]{
                    new float[] {contrast,0,0,0,0},
                    new float[] {0, contrast, 0,0,0},
                    new float[] {0,0, contrast, 0,0},
                    new float[] {0,0,0,1,0},
                    new float[] {brightness, brightness, brightness, 0,1}});
        }

        private void saveButton_Click(object sender, EventArgs e)
        {
            // TODO: SAve sub image in the proper category and all.
        }
    }
}
